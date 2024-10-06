package services

import (
	"fmt"
    "sync"
    "time"
    "os"
    "strconv"
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
}

func (bo *BatchOrchestrator) startProcessing() {
    bo.processingTicker = time.NewTicker(bo.batchDuration)
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