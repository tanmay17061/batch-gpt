package services

import (
	"batch-gpt/server/models"
	// "batch-gpt/server/logger"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

func ProcessBatch(batchRequest models.BatchRequest) (openai.ChatCompletionResponse, error) {
	request := batchRequest.Requests[0]

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	batchChatRequest := openai.CreateBatchWithUploadFileRequest{
		Endpoint:         openai.BatchEndpointChatCompletions,
		CompletionWindow: "24h",
		UploadBatchFileRequest: openai.UploadBatchFileRequest{
			FileName: "batch_request.jsonl",
			Lines: []openai.BatchLineItem{
				openai.BatchChatCompletionRequest{
					CustomID: "request_1",
					Body:     request,
					Method:   "POST",
					URL:      openai.BatchEndpointChatCompletions,
				},
			},
		},
	}

	batchResponse, err := client.CreateBatchWithUploadFile(context.Background(), batchChatRequest)
	if err != nil {
		return openai.ChatCompletionResponse{}, err
	}

	return PollAndCollectBatchResponse(client, batchResponse.ID)
}

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

func PollAndCollectBatchResponse(client *openai.Client, batchID string) (openai.ChatCompletionResponse, error) {
    ctx := context.Background()
    maxRetries := 30
    retryInterval := 2 * time.Second

    for i := 0; i < maxRetries; i++ {
        batchStatus, err := client.RetrieveBatch(ctx, batchID)
        if err != nil {
            return openai.ChatCompletionResponse{}, err
        }

        if batchStatus.Status == "completed" {
            if batchStatus.OutputFileID == nil {
                return openai.ChatCompletionResponse{}, errors.New("output file ID is missing")
            }

            rawResponse, err := client.GetFileContent(ctx, *batchStatus.OutputFileID)
            if err != nil {
                return openai.ChatCompletionResponse{}, err
            }
            defer rawResponse.Close()

            content, err := io.ReadAll(rawResponse)
            if err != nil {
                return openai.ChatCompletionResponse{}, err
            }

            lines := bytes.Split(content, []byte("\n"))
            if len(lines) == 0 {
                return openai.ChatCompletionResponse{}, errors.New("empty response file")
            }

            var batchResponseItem BatchResponseItem
            if err := json.Unmarshal(lines[0], &batchResponseItem); err != nil {
                return openai.ChatCompletionResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
            }

            if batchResponseItem.Response.StatusCode != 200 {
                return openai.ChatCompletionResponse{}, fmt.Errorf("API error: status code %d", batchResponseItem.Response.StatusCode)
            }

            if batchResponseItem.Response.Error != nil {
                return openai.ChatCompletionResponse{}, fmt.Errorf("API error: %v", batchResponseItem.Response.Error)
            }

            return batchResponseItem.Response.Body, nil
        }

        if batchStatus.Status == "failed" || batchStatus.Status == "cancelled" {
            return openai.ChatCompletionResponse{}, errors.New("batch processing failed or was cancelled")
        }

        time.Sleep(retryInterval)
    }

    return openai.ChatCompletionResponse{}, errors.New("max retries reached, batch processing did not complete in time")
}
