package cache

import (
    "batch-gpt/server/models"
    openai "github.com/sashabaranov/go-openai"
)

type Orchestrator interface {
    GetFromCache(request openai.ChatCompletionRequest) (openai.ChatCompletionResponse, bool)
    CacheResponses(requests []models.BatchRequestItem, responses []models.BatchResponseItem)
}
