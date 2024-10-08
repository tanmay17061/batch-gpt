package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"batch-gpt/server/handlers"
	"batch-gpt/server/services"
	"batch-gpt/server/db"
)

func main() {
	r := gin.Default()

	db.InitMongoDB()
	services.InitBatchOrchestrator()

	r.POST("/v1/chat/completions", handlers.HandleChatCompletions)
	r.GET("/v1/batches/:batch_id", handlers.HandleRetrieveBatch)

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}