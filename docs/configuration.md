---
layout: default
title: Configuration
nav_order: 7
---

# Configuration
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Environment Variables

### Core Settings

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `OPENAI_API_KEY` | OpenAI API key | - | Yes |
| `CLIENT_SERVING_MODE` | Service mode (sync/async/cache) | sync | No |
| `COLLATE_BATCHES_FOR_DURATION_IN_MS` | Batch window (ms) | 5000 | No |

### MongoDB Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `MONGO_HOST` | MongoDB host | localhost | No |
| `MONGO_PORT` | MongoDB port | 27017 | No |
| `MONGO_USER` | MongoDB username | admin | No |
| `MONGO_PASSWORD` | MongoDB password | password | No |
| `MONGO_DATABASE` | MongoDB database name | batchgpt | No |

### Advanced Settings

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `COLLECT_BATCH_STATS_POLLING_MAX_INTERVAL_SECONDS` | Max polling interval | 300 | No |

## Docker Compose Configuration

```yaml
version: '3.8'
services:
  mongodb:
    image: mongo:latest
    environment:
      - MONGO_INITDB_ROOT_USERNAME=admin
      - MONGO_INITDB_ROOT_PASSWORD=password
    volumes:
      - ./mongodb_data:/data/db
    ports:
      - "27017:27017"
```

## Example Configurations

### Development Setup
```bash
export OPENAI_API_KEY=your_key_here
export CLIENT_SERVING_MODE=sync
export COLLATE_BATCHES_FOR_DURATION_IN_MS=5000
export MONGO_HOST=localhost
```

### Production Setup
```bash
export OPENAI_API_KEY=your_key_here
export CLIENT_SERVING_MODE=async
export COLLATE_BATCHES_FOR_DURATION_IN_MS=3000
export COLLECT_BATCH_STATS_POLLING_MAX_INTERVAL_SECONDS=600
```

### Cache-Only Mode
```bash
export CLIENT_SERVING_MODE=cache
export MONGO_HOST=production_mongodb
```
