package models

import (
	openai "github.com/sashabaranov/go-openai"
)

// BatchResponseItem represents the structure of each item in the batch response
type BatchResponseItem struct {
    ID       string `json:"id"`
    CustomID string `json:"custom_id"`
    Response struct {
        StatusCode int                           `json:"status_code"`
        RequestID  string                        `json:"request_id"`
        Body       openai.ChatCompletionResponse `json:"body"`
        Error      *openai.APIError              `json:"error"`
    } `json:"response"`
}
