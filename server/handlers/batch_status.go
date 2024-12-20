package handlers

import (
	"batch-gpt/server/db"
	"batch-gpt/server/logger"
	"batch-gpt/services/client"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleRetrieveBatch(c *gin.Context) {
    batchID := strings.TrimPrefix(c.Param("batch_id"), "/")
    if batchID == "" {
    c.JSON(http.StatusBadRequest, openai.ErrorResponse{
	        Error: &openai.APIError{
	            Type: "invalid_request_error",
	            Message: "Missing batch_id parameter",
	        },
	    })
	    return
    }

    batchStatus, err := db.GetLatestBatchStatus(batchID)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, openai.ErrorResponse{
                Error: &openai.APIError{
                    Type: "invalid_request_error",
                    Message: "No such batch",
                },
            })
        } else {
            c.JSON(http.StatusInternalServerError, openai.ErrorResponse{
                Error: &openai.APIError{
                    Type: "internal_server_error",
                    Message: "Failed to retrieve batch status",
                },
            })
        }
        return
    }
    // Convert the Batch to BatchResponse
    response := openai.BatchResponse{
        Batch: batchStatus,
    }

    c.JSON(http.StatusOK, response)
}

func HandleListBatches(c *gin.Context) {
    batchStatuses, err := db.GetAllBatchStatuses()
    if err != nil {
        c.JSON(http.StatusInternalServerError, openai.ErrorResponse{
            Error: &openai.APIError{
                Type:    "internal_server_error",
                Message: "Failed to retrieve batch statuses",
            },
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "data": batchStatuses,
    })
}

func HandleCancelBatch(c *gin.Context) {
    batchID := strings.TrimPrefix(c.Param("batch_id"), "/")
    if batchID == "" {
        c.JSON(http.StatusBadRequest, openai.ErrorResponse{
            Error: &openai.APIError{
                Type: "invalid_request_error",
                Message: "Missing batch_id parameter",
            },
        })
        return
    }

    // Get current batch status first
    batchStatus, err := db.GetLatestBatchStatus(batchID)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, openai.ErrorResponse{
                Error: &openai.APIError{
                    Type: "invalid_request_error",
                    Message: "No such batch",
                },
            })
        } else {
            c.JSON(http.StatusInternalServerError, openai.ErrorResponse{
                Error: &openai.APIError{
                    Type: "internal_server_error",
                    Message: "Failed to retrieve batch status",
                },
            })
        }
        return
    }

    // Check if batch can be cancelled
    if batchStatus.Status == "completed" ||
       batchStatus.Status == "cancelled" ||
       batchStatus.Status == "failed" ||
       batchStatus.Status == "expired" {
        c.JSON(http.StatusBadRequest, openai.ErrorResponse{
            Error: &openai.APIError{
                Type: "invalid_request_error",
                Message: fmt.Sprintf("Cannot cancel batch in %s state", batchStatus.Status),
            },
        })
        return
    }

    // Get openAI client from context
    openAIClient := c.MustGet("openAIClient").(client.OpenAIClient)

    // Forward cancel request to OpenAI
    response, err := openAIClient.CancelBatch(c.Request.Context(), batchID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, openai.ErrorResponse{
            Error: &openai.APIError{
                Type: "internal_server_error",
                Message: "Failed to cancel batch",
            },
        })
        return
    }

    // Log the cancelled status
    err = db.LogBatchStatus(response)
    if err != nil {
        logger.WarnLogger.Printf("Failed to log cancelled batch status: %v", err)
    }

    c.JSON(http.StatusOK, response)
}
