package services

import (
	"batch-gpt/server/models"
	openai "github.com/sashabaranov/go-openai"
)

func ProcessBatch(batchRequest models.BatchRequest) (openai.ChatCompletionResponse, error) {
	// For now, we're just processing the first request in the batch
	// In the future, this is where we'd implement the actual batching logic
	request := batchRequest.Requests[0]

	// Mock response
	mockResponse := openai.ChatCompletionResponse{
		ID:      "mock_id",
		Object:  "chat.completion",
		Created: 1234567890,
		Model:   request.Model,
		Choices: []openai.ChatCompletionChoice{
			{
				Index: 0,
				Message: openai.ChatCompletionMessage{
					Role:    "assistant",
					Content: "This is a mock response",
				},
				FinishReason: "stop",
			},
		},
		Usage: openai.Usage{
			PromptTokens:     10,
			CompletionTokens: 10,
			TotalTokens:      20,
		},
	}

	return mockResponse, nil
}
