CREATE TABLE IF NOT EXISTS batch_status_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    batch_id TEXT NOT NULL,
    status TEXT NOT NULL,
    input_file_id TEXT NOT NULL,
    output_file_id TEXT,
    total_requests INTEGER DEFAULT 0,
    completed_requests INTEGER DEFAULT 0,
    failed_requests INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_batch_id ON batch_status_log(batch_id);
CREATE INDEX IF NOT EXISTS idx_created_at ON batch_status_log(created_at);