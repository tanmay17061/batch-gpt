package batch

import (
	"batch-gpt/server/models"

	openai "github.com/sashabaranov/go-openai"
)

type BatchResult struct {
    Response openai.ChatCompletionResponse
    Error    error
    IsAsync  bool
}

type Orchestrator interface {
    AddRequest(request openai.ChatCompletionRequest) <-chan BatchResult
    ProcessBatch()
    StartProcessing()
    ContinueDanglingBatches()
}

type Processor interface {
    ProcessBatch(batchRequest models.BatchRequest) ([]models.BatchResponseItem, error)
    PollAndCollectResponses(batchID string) ([]models.BatchResponseItem, error)
}
