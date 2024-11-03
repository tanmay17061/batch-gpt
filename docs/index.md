# Batch-GPT Documentation

Batch-GPT is a jump-server that converts individual OpenAI chat completion API requests into batched requests, optimizing processing and reducing costs through OpenAI's Batch API.

## Table of Contents

1. [Introduction](#introduction)
2. [Features](#features)
3. [Getting Started](#getting-started)
4. [Architecture](#architecture)
5. [Usage Guide](#usage-guide)
6. [Advanced Features](#advanced-features)
7. [Configuration](#configuration)
8. [Monitoring](#monitoring)
9. [API Reference](#api-reference)
10. [Troubleshooting](#troubleshooting)

## Introduction

Batch-GPT enables seamless integration with OpenAI's Batch API by acting as a drop-in replacement for standard OpenAI API endpoints. It intelligently collects and batches requests, providing significant cost savings while maintaining compatibility with existing OpenAI clients.

Integration is as simple as:
```diff
from openai import OpenAI
- client = OpenAI(api_key="sk-...")
+ client = OpenAI(api_key="dummy_openai_api_key", base_url="http://batch-gpt")
```

## Features

- **Cost Optimization**
  - Up to 50% savings using OpenAI's Batch API
  - Automatic request caching for zero-cost repeat queries
  
- **Reliability & Management**
  - Persistent data storage with MongoDB
  - Recovery of interrupted batches
  - Real-time batch status monitoring
  
- **Flexible Operation Modes**
  - Synchronous mode for immediate responses
  - Asynchronous mode for high-volume scenarios
  - Cache-only mode for offline operation
  
- **Security & Integration**
  - Single OpenAI key management
  - Compatible with any OpenAI-compliant client
  - Cross-session data persistence

## Getting Started

### Prerequisites
- Go 1.23.0 or later
- Docker and Docker Compose
- MongoDB (included in docker-compose setup)

### Quick Start

1. Download pre-compiled binaries from [Releases](https://github.com/tanmay17061/batch-gpt/releases)

2. Set up MongoDB:
   ```bash
   cd local/mongo
   docker-compose up -d
   ```

3. Configure environment variables:
   ```bash
   export OPENAI_API_KEY=your_openai_api_key
   export COLLATE_BATCHES_FOR_DURATION_IN_MS=5000
   export MONGO_HOST=localhost
   export MONGO_PORT=27017
   export MONGO_USER=admin
   export MONGO_PASSWORD=password
   export MONGO_DATABASE=batchgpt
   ```

4. Start the server:
   ```bash
   ./batch-gpt
   ```

## Architecture

Batch-GPT consists of several key components:

- **Server Core**: Main server handling request routing and processing
- **Batch Orchestrator**: Manages request batching and processing
- **Cache System**: Handles response caching and retrieval
- **Monitor Tool**: Terminal-based UI for batch status monitoring
- **MongoDB Backend**: Persistent storage for batches and cache

## Usage Guide

### Basic Request Flow
1. Client sends request to Batch-GPT
2. Request is checked against cache
3. If not cached:
   - Request is added to current batch
   - Batch is processed when full or timer expires
4. Response is cached and returned

### Using with Python Client
```python
from openai import OpenAI

client = OpenAI(
    api_key="dummy_openai_api_key",
    base_url="http://localhost:8080/v1"
)

response = client.chat.completions.create(
    model="gpt-3.5-turbo",
    messages=[{"role": "user", "content": "Hello!"}]
)
```

## Advanced Features

### Serving Modes

1. **Synchronous Mode** (Default)
   ```bash
   export CLIENT_SERVING_MODE=sync
   ```
   - Blocks until response is available
   - Suitable for low-volume scenarios

2. **Asynchronous Mode**
   ```bash
   export CLIENT_SERVING_MODE=async
   ```
   - Returns immediately with submission confirmation
   - Ideal for high-volume applications

3. **Cache-Only Mode**
   ```bash
   export CLIENT_SERVING_MODE=cache
   ```
   - Serves only cached responses
   - Processes pending batches from previous sessions

### Batch Monitor

Launch the monitoring tool:
```bash
./batch-monitor
```

Features:
- Real-time batch status updates
- Interactive navigation
- Status filtering
- Progress tracking

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `OPENAI_API_KEY` | OpenAI API key | Required |
| `CLIENT_SERVING_MODE` | Service mode (sync/async/cache) | sync |
| `COLLATE_BATCHES_FOR_DURATION_IN_MS` | Batch collection window | 5000 |
| `COLLECT_BATCH_STATS_POLLING_MAX_INTERVAL_SECONDS` | Max polling interval | 300 |
| `MONGO_*` | MongoDB configuration | See quick start |

## Monitoring

### Status Endpoints

- GET `/v1/batches/{batch_id}`: Get specific batch status
- GET `/v1/batches`: List all batches

### Monitor UI

The terminal-based monitor provides:
- Batch status overview
- Request counts
- Processing progress
- Error tracking

## API Reference

### Main Endpoints

- POST `/v1/chat/completions`: Chat completion requests
- GET `/v1/batches/{batch_id}`: Batch status
- GET `/v1/batches`: List batches

All endpoints maintain OpenAI API compatibility.

## Troubleshooting

### Common Issues

1. **Long Response Times**
   - Expected with OpenAI's Batch API (24h SLA)
   - Consider sync/async mode based on needs

2. **MongoDB Connection**
   - Verify MongoDB is running
   - Check connection settings

3. **Missing Responses**
   - Check batch status with monitor
   - Verify cache configuration

### Support

For issues and questions:
1. Check existing GitHub issues
2. Review the troubleshooting guide
3. Open a new issue with detailed information

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for:
- Development setup
- Code style guide
- Pull request process
- Project structure

## License

Batch-GPT is licensed under the [Apache License 2.0](LICENSE).