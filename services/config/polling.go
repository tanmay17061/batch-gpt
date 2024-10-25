package config

import (
    "batch-gpt/server/logger"
    "os"
    "time"
)

type PollingConfig interface {
    GetMaxRetryInterval() time.Duration
}

type pollingConfig struct {
    maxRetryInterval time.Duration
}

func NewPollingConfig() PollingConfig {
    maxInterval, err := time.ParseDuration(os.Getenv("COLLECT_BATCH_STATS_POLLING_MAX_INTERVAL_SECONDS") + "s")
    if err != nil {
        logger.WarnLogger.Printf("Failed to parse COLLECT_BATCH_STATS_POLLING_MAX_INTERVAL_SECONDS, using default of 300s: %v", err)
        maxInterval = 300 * time.Second
    }
    return &pollingConfig{
        maxRetryInterval: maxInterval,
    }
}

func (pc *pollingConfig) GetMaxRetryInterval() time.Duration {
    return pc.maxRetryInterval
}