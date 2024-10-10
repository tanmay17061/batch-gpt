package services

import (
    "batch-gpt/server/models"
    "bufio"
    "context"
    "encoding/json"
    "fmt"

    openai "github.com/sashabaranov/go-openai"
)

func GetBatchInputRequests(client *openai.Client, inputFileID string) ([]models.BatchRequestItem, error) {
    ctx := context.Background()

    content, err := client.GetFileContent(ctx, inputFileID)
    if err != nil {
        return nil, fmt.Errorf("failed to get file content: %w", err)
    }
    defer content.Close()

    var items []models.BatchRequestItem

    scanner := bufio.NewScanner(content)
    for scanner.Scan() {
        var batchItem struct {
            CustomID string                        `json:"custom_id"`
            Body     openai.ChatCompletionRequest `json:"body"`
        }
        if err := json.Unmarshal(scanner.Bytes(), &batchItem); err != nil {
            return nil, fmt.Errorf("failed to unmarshal batch item: %w", err)
        }
        items = append(items, models.BatchRequestItem{
            CustomID: batchItem.CustomID,
            Request:  batchItem.Body,
        })
    }

    if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("error reading file: %w", err)
    }

    return items, nil
}
