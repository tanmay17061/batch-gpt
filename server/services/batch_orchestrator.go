package services

import (
    "batch-gpt/server/db"
    "batch-gpt/server/logger"
    "batch-gpt/server/models"
    "context"
    "os"
    "strconv"
    "sync"
    "time"
    openai "github.com/sashabaranov/go-openai"
)

type BatchOrchestrator struct {
    submitNextRequests         map[string]openai.ChatCompletionRequest
    submitNextResultChannels   map[string][]chan BatchResult
    allSubmittedRequests       map[string]openai.ChatCompletionRequest
    allSubmittedResultChannels map[string][]chan BatchResult
    mu                         sync.Mutex
    batchDuration              time.Duration
    processingTicker           *time.Ticker
}

type BatchResult struct {
    Response openai.ChatCompletionResponse
    Error    error
}

var orchestrator *BatchOrchestrator

func InitBatchOrchestrator() {
    collateDuration, err := strconv.Atoi(os.Getenv("COLLATE_BATCHES_FOR_DURATION_IN_MS"))
    if err != nil {
        collateDuration = 5000 // Default to 5 seconds if not set or invalid
    }
    
    logger.InfoLogger.Printf("InitBatchOrchestrator: Batch collate duration set to %d milliseconds", collateDuration)

    orchestrator = &BatchOrchestrator{
        batchDuration:              time.Duration(collateDuration) * time.Millisecond,
        submitNextRequests:         make(map[string]openai.ChatCompletionRequest),
        submitNextResultChannels:   make(map[string][]chan BatchResult),
        allSubmittedRequests:       make(map[string]openai.ChatCompletionRequest),
        allSubmittedResultChannels: make(map[string][]chan BatchResult),
    }

    go orchestrator.startProcessing()
    go BackgroundContinueDanglingBatches()
}

func (bo *BatchOrchestrator) startProcessing() {
    bo.processingTicker = time.NewTicker(bo.batchDuration)
    
    // This loop runs indefinitely, processing batches at regular intervals.
    // The ticker fires every batchDuration, regardless of whether there are requests to process.
    // If there are no requests when the ticker fires, processBatch() will return early.
    // This approach ensures consistent batch processing intervals but may lead to some
    // unnecessary wake-ups when there are no requests to process.
    for range bo.processingTicker.C {
        bo.processBatch()
    }
}

func (bo *BatchOrchestrator) AddRequest(request openai.ChatCompletionRequest) <-chan BatchResult {
    bo.mu.Lock()
    defer bo.mu.Unlock()

    hash, err := generateRequestHash(request)
    if err != nil {
        logger.ErrorLogger.Printf("Failed to generate request hash: %v", err)
        resultChan := make(chan BatchResult, 1)
        resultChan <- BatchResult{Error: err}
        close(resultChan)
        return resultChan
    }

    resultChan := make(chan BatchResult, 1)
    
    if _, found := bo.allSubmittedRequests[hash]; found {
        logger.InfoLogger.Printf("Duplicate request found with hash: %s", hash)
        bo.allSubmittedResultChannels[hash] = append(bo.allSubmittedResultChannels[hash], resultChan)
    } else {
        logger.InfoLogger.Printf("New request added with hash: %s", hash)
        bo.submitNextRequests[hash] = request
        bo.allSubmittedRequests[hash] = request
        bo.allSubmittedResultChannels[hash] = []chan BatchResult{resultChan}
    }

    return resultChan
}

func (bo *BatchOrchestrator) processBatch() {
    bo.mu.Lock()
    requests := bo.submitNextRequests
    bo.submitNextRequests = make(map[string]openai.ChatCompletionRequest)
    bo.mu.Unlock()

    if len(requests) == 0 {
        return
    }

    batchRequest := models.BatchRequest{
        Requests: make([]models.BatchRequestItem, 0, len(requests)),
    }

    for hash, req := range requests {
        batchRequest.Requests = append(batchRequest.Requests, models.BatchRequestItem{
            CustomID: hash,
            Request:  req,
        })
    }

    responses, err := ProcessBatch(batchRequest)

    if err == nil {
        GetCacheOrchestrator().CacheResponses(batchRequest.Requests, responses)
    }

    bo.mu.Lock()
    defer bo.mu.Unlock()

    for _, response := range responses {
        result := BatchResult{
            Response: response.Response.Body,
            Error:    err,
        }
        hash := response.CustomID
        if channels, ok := bo.allSubmittedResultChannels[hash]; ok {
            for _, ch := range channels {
                ch <- result
                close(ch)
            }
            delete(bo.allSubmittedRequests, hash)
            delete(bo.allSubmittedResultChannels, hash)
        }
    }
}

func AddRequestToBatch(request openai.ChatCompletionRequest) <-chan BatchResult {
    return orchestrator.AddRequest(request)
}

func BackgroundContinueDanglingBatches() {
    logger.InfoLogger.Println("Starting to process dangling batches")
    danglingBatches, err := db.GetDanglingBatches()
    if err != nil {
        logger.ErrorLogger.Printf("Failed to get dangling batches: %v", err)
        return
    }

    logger.InfoLogger.Printf("Found %d dangling batches", len(danglingBatches))

    client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

    for _, batchID := range danglingBatches {
        go func(id string) {
            logger.InfoLogger.Printf("Processing dangling batch: %s", id)
            
            batchStatus, err := client.RetrieveBatch(context.Background(), id)
            if err != nil {
                logger.ErrorLogger.Printf("Failed to retrieve batch %s: %v", id, err)
                return
            }

            requests, err := GetBatchInputRequests(client, batchStatus.InputFileID)
            if err != nil {
                logger.ErrorLogger.Printf("Failed to get input requests for batch %s: %v", id, err)
                return
            }

            // Add dangling requests to the BatchOrchestrator
            orchestrator.mu.Lock()
            for _, req := range requests {
                hash, err := generateRequestHash(req.Request)
                if err != nil {
                    logger.ErrorLogger.Printf("Failed to generate hash for request in batch %s: %v", id, err)
                    continue
                }

                if _, exists := orchestrator.allSubmittedRequests[hash]; !exists {
                    orchestrator.allSubmittedRequests[hash] = req.Request
                    orchestrator.allSubmittedResultChannels[hash] = []chan BatchResult{}
                    logger.InfoLogger.Printf("Added dangling request with hash %s to BatchOrchestrator", hash)
                }
            }
            orchestrator.mu.Unlock()

            responses, err := PollAndCollectBatchResponses(client, id)
            if err != nil {
                logger.ErrorLogger.Printf("Failed to process dangling batch %s: %v", id, err)
                return
            }
            logger.InfoLogger.Printf("Successfully processed dangling batch: %s", id)

            // Update BatchOrchestrator with results
            orchestrator.mu.Lock()
            for _, resp := range responses {
                hash := resp.CustomID // Assuming CustomID is now the hash
                result := BatchResult{
                    Response: resp.Response.Body,
                    Error:    nil,
                }

                if channels, exists := orchestrator.allSubmittedResultChannels[hash]; exists {
                    for _, ch := range channels {
                        ch <- result
                        close(ch)
                    }
                    delete(orchestrator.allSubmittedRequests, hash)
                    delete(orchestrator.allSubmittedResultChannels, hash)
                }
            }
            orchestrator.mu.Unlock()

            // Cache the responses
            cacheRequests := make([]models.BatchRequestItem, len(requests))
            for i, req := range requests {
                hash, _ := generateRequestHash(req.Request)
                cacheRequests[i] = models.BatchRequestItem{
                    CustomID: hash,
                    Request:  req.Request,
                }
            }
            GetCacheOrchestrator().CacheResponses(cacheRequests, responses)
            logger.InfoLogger.Printf("Cached responses for dangling batch: %s", id)

            // Update batch status in the database
            err = db.LogBatchStatus(batchStatus)
            if err != nil {
                logger.ErrorLogger.Printf("Failed to update batch status for %s: %v", id, err)
            }
        }(batchID)
    }
}
