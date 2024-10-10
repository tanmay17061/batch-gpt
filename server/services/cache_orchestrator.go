package services

import (
    "batch-gpt/server/db"
    "batch-gpt/server/logger"
    "batch-gpt/server/models"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    openai "github.com/sashabaranov/go-openai"
)

type CacheOrchestrator struct {}

var cacheOrchestrator *CacheOrchestrator

func InitCacheOrchestrator() {
    cacheOrchestrator = &CacheOrchestrator{}
}

func GetCacheOrchestrator() *CacheOrchestrator {
    return cacheOrchestrator
}

func (co *CacheOrchestrator) GetFromCache(request openai.ChatCompletionRequest) (openai.ChatCompletionResponse, bool) {
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

func (co *CacheOrchestrator) CacheResponses(requests []models.BatchRequestItem, responses []models.BatchResponseItem) {
    requestMap := make(map[string]openai.ChatCompletionRequest)
    for _, req := range requests {
        requestMap[req.CustomID] = req.Request
    }

    for _, resp := range responses {
        request, ok := requestMap[resp.CustomID]
        if !ok {
            logger.WarnLogger.Printf("No matching request found for response with CustomID: %s", resp.CustomID)
            continue
        }

        hash, err := generateRequestHash(request)
        if err != nil {
            logger.ErrorLogger.Printf("Failed to generate request hash: %v", err)
            continue
        }

        err = db.CacheRequestResponse(hash, request, resp.Response.Body)
        if err != nil {
            logger.ErrorLogger.Printf("Failed to cache response: %v", err)
        } else {
            logger.InfoLogger.Printf("Cached response for request hash: %s (CustomID: %s)", hash, resp.CustomID)
        }
    }
}

func generateRequestHash(request openai.ChatCompletionRequest) (string, error) {
    requestJSON, err := json.Marshal(request)
    if err != nil {
        return "", err
    }
    hash := sha256.Sum256(requestJSON)
    return hex.EncodeToString(hash[:]), nil
}