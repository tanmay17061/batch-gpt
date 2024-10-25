package batch

import (
	"batch-gpt/server/db"
	"batch-gpt/server/logger"
	"batch-gpt/server/models"
	"batch-gpt/services/cache"
	"batch-gpt/services/config"
	"batch-gpt/services/utils"
	"context"
	// "os"
	"sync"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

type orchestrator struct {
    submitNextRequests         map[string]openai.ChatCompletionRequest
    submitNextResultChannels   map[string][]chan BatchResult
    allSubmittedRequests      map[string]openai.ChatCompletionRequest
    allSubmittedResultChannels map[string][]chan BatchResult
    mu                        sync.Mutex
    batchDuration            time.Duration
    processingTicker         *time.Ticker
    processor               Processor
    cache                   cache.Orchestrator
    servingMode             config.ServingMode
}

func NewOrchestrator(
    processor Processor,
    cache cache.Orchestrator,
    servingMode config.ServingMode,
    batchDuration time.Duration,
) *orchestrator {
    return &orchestrator{
        processor:                processor,
        cache:                    cache,
        servingMode:             servingMode,
        batchDuration:           batchDuration,
        submitNextRequests:      make(map[string]openai.ChatCompletionRequest),
        submitNextResultChannels: make(map[string][]chan BatchResult),
        allSubmittedRequests:    make(map[string]openai.ChatCompletionRequest),
        allSubmittedResultChannels: make(map[string][]chan BatchResult),
    }
}

func (bo *orchestrator) startProcessing() {
    bo.processingTicker = time.NewTicker(bo.batchDuration)

    // This loop runs indefinitely, processing batches at regular intervals.
    // The ticker fires every batchDuration, regardless of whether there are requests to process.
    // If there are no requests when the ticker fires, processBatch() will return early.
    // This approach ensures consistent batch processing intervals but may lead to some
    // unnecessary wake-ups when there are no requests to process.
    for range bo.processingTicker.C {
        go bo.processBatch()
    }
}

func (bo *orchestrator) AddRequest(request openai.ChatCompletionRequest) <-chan BatchResult {
    bo.mu.Lock()
    defer bo.mu.Unlock()

    hash, err := utils.GenerateRequestHash(request)
    if err != nil {
        logger.ErrorLogger.Printf("Failed to generate request hash: %v", err)
        resultChan := make(chan BatchResult, 1)
        resultChan <- BatchResult{Error: err}
        close(resultChan)
        return resultChan
    }

    resultChan := make(chan BatchResult, 1)

    if _, found := bo.allSubmittedRequests[hash]; found {
        logger.InfoLogger.Printf("BatchOrchestrator: cache hit: %s", hash)
    } else {
        logger.InfoLogger.Printf("BatchOrchestrator: cache miss: %s", hash)
        bo.submitNextRequests[hash] = request
        bo.allSubmittedRequests[hash] = request
        bo.allSubmittedResultChannels[hash] = []chan BatchResult{}
    }

    if bo.servingMode.IsAsync() {
        // In async mode, send an immediate result with IsAsync flag
        // and close the channel
        resultChan <- BatchResult{IsAsync: true}
        close(resultChan)
    } else {
	    // In sync mode, save channel to send result to once the response is available
	    bo.allSubmittedResultChannels[hash] = append(bo.allSubmittedResultChannels[hash], resultChan)
    }

    return resultChan
}

func (bo *orchestrator) ProcessBatch() {
    bo.processBatch()
}

func (bo *orchestrator) StartProcessing() {
    bo.startProcessing()
}

func (bo *orchestrator) processBatch() {
    bo.mu.Lock()
    requests := bo.submitNextRequests
    bo.submitNextRequests = make(map[string]openai.ChatCompletionRequest)
    bo.mu.Unlock()

    if len(requests) == 0 {
	    logger.InfoLogger.Printf("processBatch called with an empty list of requests. not submitting any batch requests.")
        return
    }

    batchRequest := models.BatchRequest{
        Requests: make([]models.BatchRequestItem, 0, len(requests)),
    }
    logger.InfoLogger.Printf("processBatch: Processing batch with %d requests", len(requests))
    for hash, req := range requests {
        batchRequest.Requests = append(batchRequest.Requests, models.BatchRequestItem{
            CustomID: hash,
            Request:  req,
        })
    }

    responses, err := bo.processor.ProcessBatch(batchRequest)

    if err == nil {
        bo.cache.CacheResponses(batchRequest.Requests, responses)
    }

    bo.mu.Lock()
    defer bo.mu.Unlock()

    for _, response := range responses {
        result := BatchResult{
            Response: response.Response.Body,
            Error:    err,
            IsAsync:  false,
        }
        hash := response.CustomID
        if channels, ok := bo.allSubmittedResultChannels[hash]; ok {
            for _, ch := range channels {
                select {
                case <-ch: // Try to receive, in case an async result was already sent
                    // Channel was already used for async response, don't send again
                default:
                    ch <- result
                }
                close(ch)
            }
            delete(bo.allSubmittedRequests, hash)
            delete(bo.allSubmittedResultChannels, hash)
        }
    }
}

func (bo *orchestrator) ContinueDanglingBatches() {
    logger.InfoLogger.Println("ContinueDanglingBatches: Starting to process dangling batches")
    danglingBatches, err := db.GetDanglingBatches()
    if err != nil {
        logger.ErrorLogger.Printf("ContinueDanglingBatches: Failed to get dangling batches: %v", err)
        return
    }

    logger.InfoLogger.Printf("ContinueDanglingBatches: Found %d dangling batches", len(danglingBatches))

    for _, batchID := range danglingBatches {
        go func(id string) {
            logger.InfoLogger.Printf("ContinueDanglingBatches: Processing dangling batch: %s", id)

            ctx := context.Background()
            batchStatus, err := bo.processor.(*processor).client.RetrieveBatch(ctx, id)
            if err != nil {
                logger.ErrorLogger.Printf("ContinueDanglingBatches: Failed to retrieve batch %s: %v", id, err)
                return
            }

            rawResponse, err := bo.processor.(*processor).client.GetFileContent(ctx, batchStatus.InputFileID)
            if err != nil {
                logger.ErrorLogger.Printf("ContinueDanglingBatches: Failed to get file content: %v", err)
                return
            }

            requests, err := GetBatchInputRequests(rawResponse)
            if err != nil {
                logger.ErrorLogger.Printf("ContinueDanglingBatches: Failed to parse input requests: %v", err)
                return
            }

            // Add dangling requests to the BatchOrchestrator
            bo.mu.Lock()
            for _, req := range requests {
                hash, err := utils.GenerateRequestHash(req.Request)
                if err != nil {
                    logger.ErrorLogger.Printf("ContinueDanglingBatches: Failed to generate hash for request in batch %s: %v", id, err)
                    continue
                }

                if _, exists := bo.allSubmittedRequests[hash]; !exists {
                    bo.allSubmittedRequests[hash] = req.Request
                    bo.allSubmittedResultChannels[hash] = []chan BatchResult{}
                    logger.InfoLogger.Printf("ContinueDanglingBatches: Added dangling request with hash %s to BatchOrchestrator", hash)
                }
            }
            bo.mu.Unlock()

            responses, err := bo.processor.PollAndCollectResponses(id)
            if err != nil {
                logger.ErrorLogger.Printf("ContinueDanglingBatches: Failed to process dangling batch %s: %v", id, err)
                return
            }
            logger.InfoLogger.Printf("ContinueDanglingBatches: Successfully processed dangling batch: %s", id)

            // Update BatchOrchestrator with results
            bo.mu.Lock()
            for _, resp := range responses {
                hash := resp.CustomID // Assuming hash was submitted as the custom id during batch request creation to openAI
                result := BatchResult{
                    Response: resp.Response.Body,
                    Error:    nil,
                    // if a new request arrives for a dangling batch in sync mode,
                    // it needs to receive IsAsync as false.
                    IsAsync:  false,
                }

                // In case of a dangling batch, orchestrator.allSubmittedResultChannels[hash] will
                // contain channels for requests that were accumulated while the dangline batch was being processed.
                if channels, exists := bo.allSubmittedResultChannels[hash]; exists {
                    logger.InfoLogger.Printf("ContinueDanglingBatches: Dangling batch with hash=%s has %d result channels", hash, len(channels))
                    for _, ch := range channels {
                        select {
                        case ch <- result:
                            // Successfully sent the result
                            logger.InfoLogger.Printf("ContinueDanglingBatches: Result successfully sent to channel")
                        default:
                            // Channel is full or closed, log this situation
                            logger.WarnLogger.Printf("Unable to send result for dangling batch item %s", hash)
                        }
                        close(ch)
                    }
                    delete(bo.allSubmittedRequests, hash)
                    delete(bo.allSubmittedResultChannels, hash)
                }
            }
            bo.mu.Unlock()

            // Cache the responses
            cacheRequests := make([]models.BatchRequestItem, len(requests))
            for i, req := range requests {
                hash, _ := utils.GenerateRequestHash(req.Request)
                cacheRequests[i] = models.BatchRequestItem{
                    CustomID: hash,
                    Request:  req.Request,
                }
            }

            bo.cache.CacheResponses(cacheRequests, responses)
            logger.InfoLogger.Printf("ContinueDanglingBatches: Cached responses for dangling batch: %s", id)

            // Update batch status in the database
            err = db.LogBatchStatus(batchStatus)
            if err != nil {
                logger.ErrorLogger.Printf("ContinueDanglingBatches: Failed to update batch status for %s: %v", id, err)
            }
        }(batchID)
    }
}
