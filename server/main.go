package main

import (
    "batch-gpt/server/db"
    "batch-gpt/server/handlers"
    "batch-gpt/services/batch"
    "batch-gpt/services/cache"
    "batch-gpt/services/client"
    "batch-gpt/services/config"
    "log"
    "os"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
)

func main() {
    // Initialize configurations
    servingMode := config.NewServingMode(os.Getenv("CLIENT_SERVING_MODE"))
    pollingConfig := config.NewPollingConfig()

    // Initialize database
    db.InitMongoDB()

    // Initialize services
    openAIClient := client.NewOpenAIClient(os.Getenv("OPENAI_API_KEY"))
    cacheOrch := cache.NewOrchestrator()
    
    // Get batch duration from env
    collateDuration, err := strconv.Atoi(os.Getenv("COLLATE_BATCHES_FOR_DURATION_IN_MS"))
    if err != nil {
        collateDuration = 5000 // Default to 5 seconds if not set or invalid
    }
    batchDuration := time.Duration(collateDuration) * time.Millisecond

    // Initialize batch processor and orchestrator
    batchProcessor := batch.NewProcessor(openAIClient, pollingConfig)
    batchOrch := batch.NewOrchestrator(
        batchProcessor,
        cacheOrch,
        servingMode,
        batchDuration,
    )

    if servingMode.IsCache() {
        // In cache mode, only process dangling batches
        go batchOrch.ContinueDanglingBatches()
        log.Println("Server starting in cache-only mode - processing only dangling batches")
    } else {
        // In sync/async mode, start regular processing and handle dangling batches
        go batchOrch.StartProcessing()
        go batchOrch.ContinueDanglingBatches()
        log.Println("Server starting in", servingMode.GetMode(), "mode - processing new and dangling batches")
    }

    // Initialize router
    r := gin.Default()

    // Update handlers to use the new services
    r.POST("/v1/chat/completions", handlers.NewChatCompletionsHandler(batchOrch, cacheOrch, servingMode))
    r.GET("/v1/batches/:batch_id", handlers.HandleRetrieveBatch)
    r.GET("/v1/batches", handlers.HandleListBatches)

    log.Println("Server starting on :8080")
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}