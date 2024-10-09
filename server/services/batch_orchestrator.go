package services

import (
	"fmt"
    "sync"
    "time"
    "os"
    "strconv"
    "batch-gpt/server/db"
    "batch-gpt/server/models"
    "batch-gpt/server/logger"
    openai "github.com/sashabaranov/go-openai"
)

type BatchOrchestrator struct {
    requests         map[string]openai.ChatCompletionRequest
    resultChannels   map[string]chan BatchResult
    mu               sync.Mutex
    batchDuration    time.Duration
    processingTicker *time.Ticker
    nextRequestID    int
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
        batchDuration:  time.Duration(collateDuration) * time.Millisecond,
        requests:       make(map[string]openai.ChatCompletionRequest),
        resultChannels: make(map[string]chan BatchResult),
        nextRequestID:  1,
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

func (bo *BatchOrchestrator) processBatch() {
    bo.mu.Lock()
    requests := bo.requests
    channels := bo.resultChannels
    bo.requests = make(map[string]openai.ChatCompletionRequest)
    bo.resultChannels = make(map[string]chan BatchResult)
    bo.mu.Unlock()

    if len(requests) == 0 {
        return
    }

    batchRequest := models.BatchRequest{
        Requests: make([]openai.ChatCompletionRequest, 0, len(requests)),
    }
    requestIDs := make([]string, 0, len(requests))

    for id, req := range requests {
        batchRequest.Requests = append(batchRequest.Requests, req)
        requestIDs = append(requestIDs, id)
    }

    responses, err := ProcessBatch(batchRequest)

    for i, id := range requestIDs {
        result := BatchResult{Error: err}
        if err == nil && i < len(responses) {
            result.Response = responses[i]
        }
        channels[id] <- result
        close(channels[id])
    }
}

func (bo *BatchOrchestrator) AddRequest(request openai.ChatCompletionRequest) <-chan BatchResult {
    bo.mu.Lock()
    defer bo.mu.Unlock()

    id := fmt.Sprintf("req_%d", bo.nextRequestID)
    
    logger.InfoLogger.Printf("BatchOrchestrator.AddRequest: Request %s added to the batch", id)

    bo.nextRequestID++

    bo.requests[id] = request
    resultChan := make(chan BatchResult, 1)
    bo.resultChannels[id] = resultChan

    return resultChan
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

    for _, batchID := range danglingBatches {
        go func(id string) {
            logger.InfoLogger.Printf("Processing dangling batch: %s", id)
            client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
            responses, err := PollAndCollectBatchResponses(client, id)
            if err != nil {
                logger.ErrorLogger.Printf("Failed to process dangling batch %s: %v", id, err)
                return
            }
            logger.InfoLogger.Printf("Successfully processed dangling batch: %s", id)
            err = db.CacheResponses(id, responses)
            if err != nil {
                logger.ErrorLogger.Printf("Failed to cache responses for batch %s: %v", id, err)
            } else {
                logger.InfoLogger.Printf("Cached responses for batch: %s", id)
            }
        }(batchID)
    }
}