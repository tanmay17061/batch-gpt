package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"batch-gpt/server/handlers"
)

func main() {
	r := gin.Default()

	r.POST("/v1/chat/completions", handlers.HandleChatCompletions)

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}