package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    openai "github.com/sashabaranov/go-openai"
    "batch-gpt/server/services"
)

func HandleChatCompletions(c *gin.Context) {
    var request openai.ChatCompletionRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check cache first
    if cachedResponse, found := services.GetCacheOrchestrator().GetFromCache(request); found {
        c.JSON(http.StatusOK, cachedResponse)
        return
    }

    // If not in cache, proceed with batch processing
    resultChan := services.AddRequestToBatch(request)

    select {
    case result := <-resultChan:
        if result.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        } else {
            c.JSON(http.StatusOK, result.Response)
        }
    case <-c.Request.Context().Done():
        c.JSON(http.StatusRequestTimeout, gin.H{"error": "Request timeout"})
    }
}
