package handlers

import (
	"batch-gpt/server/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGetBatchStatus(c *gin.Context) {
    batchID := c.Query("batch_id")
    if batchID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "batch_id is required"})
        return
    }

    status, counts, err := db.GetLatestBatchStatus(batchID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve batch status"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "batch_id": batchID,
        "status": status,
        "request_counts": counts,
    })
}
