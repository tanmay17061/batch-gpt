package client

import (
    "context"
    openai "github.com/sashabaranov/go-openai"
)

type openAIClient struct {
    client *openai.Client
}

func NewOpenAIClient(apiKey string) OpenAIClient {
    return &openAIClient{
        client: openai.NewClient(apiKey),
    }
}

func (c *openAIClient) CreateBatchWithUploadFile(ctx context.Context, req openai.CreateBatchWithUploadFileRequest) (openai.BatchResponse, error) {
    return c.client.CreateBatchWithUploadFile(ctx, req)
}

func (c *openAIClient) RetrieveBatch(ctx context.Context, batchID string) (openai.BatchResponse, error) {
    return c.client.RetrieveBatch(ctx, batchID)
}

func (c *openAIClient) GetFileContent(ctx context.Context, fileID string) (openai.RawResponse, error) {
    return c.client.GetFileContent(ctx, fileID)
}
