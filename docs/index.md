# Batch-GPT Documentation

Batch-GPT is an efficient service that converts individual OpenAI chat completion API requests into batched requests, optimizing processing and potentially reducing costs for high-volume applications.

## Table of Contents

1. [Introduction](#introduction)
2. [Features](#features)
3. [Installation](#installation)
4. [Usage](#usage)
5. [API Reference](#api-reference)
6. [Configuration](#configuration)
7. [Performance](#performance)
8. [Troubleshooting](#troubleshooting)
9. [Contributing](#contributing)
10. [License](#license)

## Introduction

Batch-GPT acts as a middleware between your application and OpenAI's API. It collects individual chat completion requests over a short time window, bundles them into a single batch request, and then sends this batch to OpenAI. This approach can lead to significant improvements in throughput and potential cost savings for applications that generate a high volume of requests.

## Features

- Automatic batching of chat completion requests
- Compatible with OpenAI's API structure
- Configurable batch window duration
- Efficient handling of large volumes of requests
- Potential for reduced API costs
- Easy integration with existing applications

## Installation

To install Batch-GPT, follow these steps:

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/batch-gpt.git
   ```
2. Navigate to the project directory:
   ```
   cd batch-gpt
   ```
3. Install dependencies:
   ```
   go mod tidy
   ```
4. Set up the SQLite database:
   ```
   ./setup_db.sh
   ```

## Usage

To start the Batch-GPT server:

```
go run server/main.go
```

The server will start on `localhost:8080`. You can now send requests to this server as if you were interacting with the OpenAI API directly.

Example using curl:

```bash
curl http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

## API Reference

Batch-GPT mimics OpenAI's API endpoints. The main endpoint is:

- POST `/v1/chat/completions`: Accept chat completion requests

Additional endpoints:

- GET `/get-batch-status`: Retrieve the status of a specific batch

For detailed API documentation, refer to [OpenAI's API Reference](https://platform.openai.com/docs/api-reference/chat).

## Configuration

Batch-GPT can be configured using environment variables:

- `OPENAI_API_KEY`: Your OpenAI API key
- `COLLATE_BATCHES_FOR_DURATION_IN_MS`: The duration to collect requests before sending a batch (default: 5000ms)

## Performance

Batch-GPT is designed to handle high volumes of requests efficiently. In testing, it has shown significant improvements in throughput compared to individual API calls, especially for applications generating 50+ requests per second.

## Troubleshooting

Common issues and their solutions:

1. **Connection refused**: Ensure the server is running and you're connecting to the correct port.
2. **API key errors**: Verify that your OpenAI API key is correctly set in the environment variables.
3. **Batch processing delays**: If responses seem delayed, check the `COLLATE_BATCHES_FOR_DURATION_IN_MS` setting. A lower value will result in more frequent, smaller batches.

## Contributing

We welcome contributions to Batch-GPT! Please see our [Contributing Guide](CONTRIBUTING.md) for more details on how to get started.

## License

Batch-GPT is released under the [Apache License 2.0](LICENSE).