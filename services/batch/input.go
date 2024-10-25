package batch

import (
    "batch-gpt/server/models"
    "bufio"
    "encoding/json"
    "fmt"
    "io"
    openai "github.com/sashabaranov/go-openai"
)

func GetBatchInputRequests(rawResponse io.ReadCloser) ([]models.BatchRequestItem, error) {
    defer rawResponse.Close()

    var items []models.BatchRequestItem
    scanner := bufio.NewScanner(rawResponse)

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
