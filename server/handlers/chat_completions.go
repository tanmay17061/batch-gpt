package handlers

import (
    "batch-gpt/services/batch"
    "batch-gpt/services/cache"
    "batch-gpt/services/config"
    "net/http"
    "github.com/gin-gonic/gin"
    openai "github.com/sashabaranov/go-openai"
)

type ChatCompletionsHandler struct {
    batchOrch batch.Orchestrator
    cacheOrch cache.Orchestrator
    servingMode config.ServingMode
}

func NewChatCompletionsHandler(batchOrch batch.Orchestrator, cacheOrch cache.Orchestrator, servingMode config.ServingMode) gin.HandlerFunc {
    handler := &ChatCompletionsHandler{
        batchOrch: batchOrch,
        cacheOrch: cacheOrch,
        servingMode: servingMode,
    }
    return handler.Handle
}

func (h *ChatCompletionsHandler) Handle(c *gin.Context) {
    var request openai.ChatCompletionRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check cache first
    if cachedResponse, found := h.cacheOrch.GetFromCache(request); found {
        c.JSON(http.StatusOK, cachedResponse)
        return
    }

    // If cache-only mode and no cache hit, return error
    if h.servingMode.IsCache() {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Response not found in cache and server is in cache-only mode",
        })
        return
    }

    // Normal processing for async/sync modes
    resultChan := h.batchOrch.AddRequest(request)
    
    select {
    case result := <-resultChan:
        if result.IsAsync {
            c.JSON(http.StatusAccepted, gin.H{"message": "Request submitted for processing"})
            return
        }
        if result.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        } else {
            c.JSON(http.StatusOK, result.Response)
        }
    case <-c.Request.Context().Done():
        c.JSON(http.StatusRequestTimeout, gin.H{"error": "Request timeout"})
    }
}