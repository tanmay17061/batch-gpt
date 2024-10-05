package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
)

func main() {
	r := gin.Default()

	r.POST("/v1/chat/completions", handleChatCompletions)

	r.Run(":8080")
}

func handleChatCompletions(c *gin.Context) {
	var request openai.ChatCompletionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// For now, we'll just echo back the request to confirm it's working
	c.JSON(http.StatusOK, request)
}