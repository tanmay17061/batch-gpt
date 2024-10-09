package services

import (
	"batch-gpt/server/db"
	"batch-gpt/server/logger"
	"batch-gpt/server/models"
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

func ProcessBatch(batchRequest models.BatchRequest) ([]openai.ChatCompletionResponse, error) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	batchChatRequest := openai.CreateBatchWithUploadFileRequest{
		Endpoint:         openai.BatchEndpointChatCompletions,
		CompletionWindow: "24h",
		UploadBatchFileRequest: openai.UploadBatchFileRequest{
			FileName: "batch_request.jsonl",
			Lines:    make([]openai.BatchLineItem, len(batchRequest.Requests)),
		},
	}

	for i, request := range batchRequest.Requests {
		batchChatRequest.UploadBatchFileRequest.Lines[i] = openai.BatchChatCompletionRequest{
			CustomID: fmt.Sprintf("request_%d", i+1),
			Body:     request,
			Method:   "POST",
			URL:      openai.BatchEndpointChatCompletions,
		}
	}

	batchResponse, err := client.CreateBatchWithUploadFile(context.Background(), batchChatRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create batch: %w", err)
	}

	// Log initial batch status
	err = db.LogBatchResponse(batchResponse)
    if err != nil {
        logger.WarnLogger.Printf("Failed to log initial batch status: %v", err)
    }

	responses, err := PollAndCollectBatchResponses(client, batchResponse.ID)
	if err != nil {
        logger.ErrorLogger.Printf("Failed to process live batch %s: %v", batchResponse.ID, err)
        return responses, err
    }

    logger.InfoLogger.Printf("Successfully processed live batch: %s", batchResponse.ID)
    
    go func() {
        err := db.CacheResponses(batchResponse.ID, responses)
        if err != nil {
            logger.ErrorLogger.Printf("Failed to cache responses for live batch %s: %v", batchResponse.ID, err)
        } else {
            logger.InfoLogger.Printf("Cached responses for live batch: %s", batchResponse.ID)
        }
    }()

    return responses, err
}

func PollAndCollectBatchResponses(client *openai.Client, batchID string) ([]openai.ChatCompletionResponse, error) {
	ctx := context.Background()
	maxRetries := 60
	retryInterval := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		batchStatus, err := client.RetrieveBatch(ctx, batchID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve batch status: %w", err)
		}

		logger.InfoLogger.Printf("Batch Status: ID=%s, Status=%s, InputFileID=%s, OutputFileID=%v, RequestCounts=%+v",
			batchStatus.ID, batchStatus.Status, batchStatus.InputFileID, batchStatus.OutputFileID, batchStatus.RequestCounts)

		// Log current batch status
		err = db.LogBatchResponse(batchStatus)
        if err != nil {
            logger.WarnLogger.Printf("Failed to log batch status: %v", err)
        }

		if batchStatus.Status == "completed" {
			if batchStatus.OutputFileID == nil {
				return nil, errors.New("output file ID is missing")
			}

			rawResponse, err := client.GetFileContent(ctx, *batchStatus.OutputFileID)
			if err != nil {
				return nil, fmt.Errorf("failed to get file content: %w", err)
			}
			defer rawResponse.Close()

			content, err := io.ReadAll(rawResponse)
			if err != nil {
				return nil, fmt.Errorf("failed to read response content: %w", err)
			}

			lines := bytes.Split(content, []byte("\n"))
			responses := make([]openai.ChatCompletionResponse, 0, len(lines))

			for l_i, line := range lines {
				if len(line) == 0 {
					// Skip empty lines
					continue
				}

				logger.InfoLogger.Printf("Line %d contents: %s", l_i+1, string(line))
				var batchResponseItem models.BatchResponseItem
				if err := json.Unmarshal(line, &batchResponseItem); err != nil {
					return nil, fmt.Errorf("failed to unmarshal response item: %w", err)
				}

				if batchResponseItem.Response.StatusCode != 200 {
					return nil, fmt.Errorf("API error for item %s: status code %d", batchResponseItem.CustomID, batchResponseItem.Response.StatusCode)
				}

				if batchResponseItem.Response.Error != nil {
					return nil, fmt.Errorf("API error for item %s: %v", batchResponseItem.CustomID, batchResponseItem.Response.Error)
				}

				responses = append(responses, batchResponseItem.Response.Body)
			}

			return responses, nil
		}

		if batchStatus.Status == "failed" || batchStatus.Status == "cancelled" {
			return nil, fmt.Errorf("batch processing %s", batchStatus.Status)
		}

		time.Sleep(retryInterval)
	}

	return nil, errors.New("max retries reached, batch processing did not complete in time")
}
