package handlers

import (
	"net/http"

	"batch-gpt/server/models"
	"batch-gpt/server/services"

	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
)

func HandleChatCompletions(c *gin.Context) {
	var request openai.ChatCompletionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	batchRequest := models.BatchRequest{
		Requests: []openai.ChatCompletionRequest{request},
	}

	response, err := services.ProcessBatch(batchRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
