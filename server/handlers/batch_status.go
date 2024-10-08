package handlers

import (
	"batch-gpt/server/db"
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
