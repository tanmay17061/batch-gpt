package client

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIClient interface {
	CreateBatchWithUploadFile(context.Context, openai.CreateBatchWithUploadFileRequest) (openai.BatchResponse, error)
	RetrieveBatch(context.Context, string) (openai.BatchResponse, error)
	GetFileContent(context.Context, string) (openai.RawResponse, error)
	CancelBatch(context.Context, string) (openai.BatchResponse, error)
}
