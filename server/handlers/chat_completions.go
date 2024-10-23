package handlers

import (
    "batch-gpt/services/batch"
    "batch-gpt/services/cache"
    "net/http"
    "github.com/gin-gonic/gin"
    openai "github.com/sashabaranov/go-openai"
)

type ChatCompletionsHandler struct {
    batchOrch batch.Orchestrator
    cacheOrch cache.Orchestrator
}

func NewChatCompletionsHandler(batchOrch batch.Orchestrator, cacheOrch cache.Orchestrator) gin.HandlerFunc {
    handler := &ChatCompletionsHandler{
        batchOrch: batchOrch,
        cacheOrch: cacheOrch,
    }
    return handler.Handle
}

func (h *ChatCompletionsHandler) Handle(c *gin.Context) {
    var request openai.ChatCompletionRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if cachedResponse, found := h.cacheOrch.GetFromCache(request); found {
        c.JSON(http.StatusOK, cachedResponse)
        return
    }

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