package models

import (
	openai "github.com/sashabaranov/go-openai"
)

type BatchRequestItem struct {
    CustomID string
    Request  openai.ChatCompletionRequest
}

type BatchRequest struct {
    Requests []BatchRequestItem
}