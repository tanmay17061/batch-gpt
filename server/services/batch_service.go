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

var maxRetryIntervalSeconds time.Duration

func InitPollingParameters() {
	maxInterval, err := time.ParseDuration(os.Getenv("COLLECT_BATCH_POLLING_MAX_INTERVAL_SECONDS") + "s")
	if err != nil {
		logger.WarnLogger.Printf("Failed to parse COLLECT_BATCH_MAX_INTERVAL_SECONDS, using default of 300s: %v", err)
		maxRetryIntervalSeconds = 300 * time.Second
	} else {
		maxRetryIntervalSeconds = maxInterval
	}
}

func ProcessBatch(batchRequest models.BatchRequest) ([]models.BatchResponseItem, error) {
    client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

    batchChatRequest := openai.CreateBatchWithUploadFileRequest{
        Endpoint:         openai.BatchEndpointChatCompletions,
        CompletionWindow: "24h",
        UploadBatchFileRequest: openai.UploadBatchFileRequest{
            FileName: "batch_request.jsonl",
            Lines:    make([]openai.BatchLineItem, len(batchRequest.Requests)),
        },
    }

    for i, requestItem := range batchRequest.Requests {
        batchChatRequest.UploadBatchFileRequest.Lines[i] = openai.BatchChatCompletionRequest{
            CustomID: requestItem.CustomID,
            Body:     requestItem.Request,
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

	return responses, err
}

func PollAndCollectBatchResponses(client *openai.Client, batchID string) ([]models.BatchResponseItem, error) {
	ctx := context.Background()

	retryIntervalSeconds := 5 * time.Second

	for {
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
			var responses []models.BatchResponseItem
	        for _, line := range lines {
	            if len(line) == 0 {
	                continue
	            }
	
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
	
	            responses = append(responses, batchResponseItem)
	        }
	
	        return responses, nil
		}

		if batchStatus.Status == "failed" || batchStatus.Status == "cancelled" {
			return nil, fmt.Errorf("batch processing %s", batchStatus.Status)
		}

		time.Sleep(retryIntervalSeconds)
		if retryIntervalSeconds < maxRetryIntervalSeconds {
			retryIntervalSeconds = retryIntervalSeconds * 2
			if retryIntervalSeconds > maxRetryIntervalSeconds {
				retryIntervalSeconds = maxRetryIntervalSeconds
			}
		}
	}
}
