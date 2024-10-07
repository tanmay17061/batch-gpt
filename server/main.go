package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"batch-gpt/server/handlers"
	"batch-gpt/server/services"
)

func main() {
	r := gin.Default()

	services.InitBatchOrchestrator()

	r.POST("/v1/chat/completions", handlers.HandleChatCompletions)
	r.GET("/get-batch-status", handlers.HandleGetBatchStatus)

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}