package db

import (
    "database/sql"
    "log"
    "sync"

    _ "github.com/mattn/go-sqlite3"
    "github.com/sashabaranov/go-openai"
)

var (
    db   *sql.DB
    once sync.Once
)

func GetDB() *sql.DB {
    once.Do(func() {
        var err error
        db, err = sql.Open("sqlite3", "./batchgpt.db")
        if err != nil {
            log.Fatalf("Failed to open database: %v", err)
        }
    })
    return db
}

func LogBatchStatus(batchID, status, inputFileID string, outputFileID *string, requestCounts openai.BatchRequestCounts) error {
    var err error
    if outputFileID == nil {
        _, err = GetDB().Exec(`
            INSERT INTO batch_status_log (
                batch_id, status, input_file_id, output_file_id,
                total_requests, completed_requests, failed_requests
            ) VALUES (?, ?, ?, NULL, ?, ?, ?)
        `, batchID, status, inputFileID, requestCounts.Total, requestCounts.Completed, requestCounts.Failed)
    } else {
        _, err = GetDB().Exec(`
            INSERT INTO batch_status_log (
                batch_id, status, input_file_id, output_file_id,
                total_requests, completed_requests, failed_requests
            ) VALUES (?, ?, ?, ?, ?, ?, ?)
        `, batchID, status, inputFileID, *outputFileID, requestCounts.Total, requestCounts.Completed, requestCounts.Failed)
    }
    return err
}

func GetLatestBatchStatus(batchID string) (string, openai.BatchRequestCounts, error) {
    var status string
    var counts openai.BatchRequestCounts
    err := GetDB().QueryRow(`
        SELECT status, total_requests, completed_requests, failed_requests 
        FROM batch_status_log 
        WHERE batch_id = ?
        ORDER BY created_at DESC
        LIMIT 1
    `, batchID).Scan(&status, &counts.Total, &counts.Completed, &counts.Failed)
    return status, counts, err
}
