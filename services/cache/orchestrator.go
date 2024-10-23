package cache

import (
	"batch-gpt/server/db"
	"batch-gpt/server/logger"
	"batch-gpt/server/models"

	openai "github.com/sashabaranov/go-openai"
)

type orchestrator struct{}

func NewOrchestrator() Orchestrator {
    return &orchestrator{}
}

func (co *orchestrator) GetFromCache(request openai.ChatCompletionRequest) (openai.ChatCompletionResponse, bool) {
    hash, err := generateRequestHash(request)
    if err != nil {
        logger.ErrorLogger.Printf("Failed to generate request hash: %v", err)
        return openai.ChatCompletionResponse{}, false
    }

    cachedResponse, err := db.GetCachedResponse(hash)
    if err == nil {
        logger.InfoLogger.Printf("Cache hit for request hash: %s", hash)
        return cachedResponse, true
    }

    logger.InfoLogger.Printf("Cache miss for request hash: %s", hash)
    return openai.ChatCompletionResponse{}, false
}

func (co *orchestrator) CacheResponses(requests []models.BatchRequestItem, responses []models.BatchResponseItem) {
    requestMap := make(map[string]openai.ChatCompletionRequest)
    for _, req := range requests {
        requestMap[req.CustomID] = req.Request
    }
    failed_caches := 0
    success_caches := 0

    for _, resp := range responses {
        request, ok := requestMap[resp.CustomID]
        if !ok {
            logger.WarnLogger.Printf("No matching request found for response with CustomID: %s", resp.CustomID)
            failed_caches += 1
            continue
        }

        hash, err := generateRequestHash(request)
        if err != nil {
            logger.ErrorLogger.Printf("Failed to generate request hash: %v", err)
            failed_caches += 1
            continue
        }

        err = db.CacheRequestResponse(hash, request, resp.Response.Body)
        if err != nil {
        	failed_caches += 1
        } else {
	        success_caches += 1
        }
    }
    logger.InfoLogger.Printf("Caching results: %d/%d successful, %d failed", success_caches, len(responses), failed_caches)
}
