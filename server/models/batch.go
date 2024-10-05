package models

import (
	openai "github.com/sashabaranov/go-openai"
)

type BatchRequest struct {
	Requests []openai.ChatCompletionRequest
}