# batch-gpt

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/tanmay17061/batch-gpt)](https://github.com/tanmay17061/batch-gpt/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/tanmay17061/batch-gpt/release.yml)](https://github.com/tanmay17061/batch-gpt/actions)
[![License](https://img.shields.io/github/license/tanmay17061/batch-gpt)](https://github.com/tanmay17061/batch-gpt/blob/main/LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/tanmay17061/batch-gpt)](https://github.com/tanmay17061/batch-gpt/blob/main/go.mod)
[![Last Commit](https://img.shields.io/github/last-commit/tanmay17061/batch-gpt)](https://github.com/tanmay17061/batch-gpt/commits/main)

A jump-server to convert openai chat completion api requests to batched chat completion requests

It really is as simple as:
```diff
from openai import OpenAI
- client = OpenAI(api_key="sk-...")
+ client = OpenAI(api_key="dummy_openai_api_key", base_url="http://batch-gpt")
```
> Just an example, Batch-GPT works with any OpenAI-compatible client.

## Features 🤓👍

- Seamless Integration: Drop-in replacement for standard OpenAI API clients
- Cost-Effective:
  1. ⭐ Up to 50% savings using OpenAI's Batch API
  2. Automatic request caching for zero-cost repeat queries
- Enhanced Reliability: Resumes processing of interrupted batches on server restart
- Persistent Data: MongoDB integration for cross-session data retention
- Centralized Management: View all batch statuses at once
- Flexible Serving Modes:
  - Synchronous mode for immediate responses
  - Asynchronous mode for handling high-volume requests
- Secure Key Distribution: Single OpenAI key for all clients, maintained via Batch-GPT

## Limitations 🤔

- High Turnaround Time: OpenAI's Batch API has a 24-hour SLA (as of 10-10-2024).
- Not Suitable for Real-Time: Potential delays make it unsuitable for live requests
- Reliability Measures: While implemented, may not fully mitigate long processing times

> **💡** Consider OpenAI's [Realtime API](https://platform.openai.com/docs/guides/realtime) for immediate response needs.

## Prerequisites

- Go 1.23.0 or later
- Docker and (Docker Compose if running MongoDB through `local/mongo/docker-compose.yaml`)
- An OpenAI API key

## Setup

You can either build the server from source (for the latest changes) or download pre-compiled binaries.

### Option 1: Using Pre-compiled Binaries

1. Download the latest release for your operating system (darwin/linux/windows) and architecture (amd64/arm64) from the [Releases page](https://github.com/tanmay17061/batch-gpt/releases).

2. Extract the downloaded archive.

3. Set up the MongoDB database:
   ```
   cd local/mongo
   docker-compose up -d
   cd ../..
   ```

4. Set environment variables:
   ```
   export OPENAI_API_KEY=your_openai_api_key_here
   export COLLATE_BATCHES_FOR_DURATION_IN_MS=5000
   export MONGO_HOST=localhost
   export MONGO_PORT=27017
   export MONGO_USER=admin
   export MONGO_PASSWORD=password
   export MONGO_DATABASE=batchgpt
   ```

5. Run the server:
   ```
   ./batch-gpt
   ```

### Option 2: Building from Source

1. Clone the repository:
   ```
   git clone https://github.com/tanmay17061/batch-gpt.git
   cd batch-gpt
   ```

2. Set up the MongoDB database:
   ```
   cd local/mongo
   docker-compose up -d
   cd ../..
   ```

3. Set environment variables:
   ```
   export OPENAI_API_KEY=your_openai_api_key_here
   export COLLATE_BATCHES_FOR_DURATION_IN_MS=5000
   export MONGO_HOST=localhost
   export MONGO_PORT=27017
   export MONGO_USER=admin
   export MONGO_PASSWORD=password
   export MONGO_DATABASE=batchgpt
   ```

4. Build and run the server:
   ```
   go build -o batch-gpt server/main.go
   ./batch-gpt
   ```

The server will start on `http://localhost:8080`.

## Usage
Note: In asynchronous mode, the server will return immediately with a submission confirmation instead of waiting for the actual response. Look at the Advanced Settings section to learn more about sync/async modes.

### Sending Chat Completion Requests

You can send requests to the batch-gpt server using any existing openai client.

#### Using curl
Send POST requests to `/v1/chat/completions` with the same format as the OpenAI API. For example:

```bash
curl http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

#### Using OpenAI Python Library

Make sure you have the OpenAI Python library installed:

```bash
pip install openai
```

```python
from openai import OpenAI

# Create a custom OpenAI client that points to the batch-gpt server
client = OpenAI(
    api_key="dummy_openai_api_key",  # The API key is not used by batch-gpt, but is required by the client
    base_url="http://localhost:8080/v1"  # Point to your batch-gpt server
)

# Send a chat completion request
chat_completion = client.chat.completions.create(
    model="gpt-3.5-turbo",
    messages=[
        {"role": "user", "content": "Hello!"}
    ]
)

# Print the response
print(chat_completion.choices[0].message.content)
```

### Checking Batch Status

You can check the status of a batch using any existing openai client.

#### Using curl

To check the status of a specific batch:

```bash
curl http://localhost:8080/v1/batches/{your_batch_id_here}
```

To retrieve the status of all batches:

```bash
curl http://localhost:8080/v1/batches
```

#### Using OpenAI Python Client

You can use the OpenAI Python client to check batch statuses. Here's an example:

```python
from openai import OpenAI

# Create a custom OpenAI client that points to the batch-gpt server
client = OpenAI(
    api_key="dummy_openai_api_key",  # The API key is not used by batch-gpt, but is required by the client
    base_url="http://localhost:8080/v1"  # Point to your batch-gpt server
)

# Retrieve the status of a specific batch
batch_id = "your_batch_id_here"
batch_status = client.batches.retrieve(batch_id)

# Print the batch status
print(f"Batch ID: {batch_status.batch.id}")
print(f"Status: {batch_status.batch.status}")
print(f"Created At: {batch_status.batch.created_at}")
print(f"Expires At: {batch_status.batch.expires_at}")
print(f"Request Counts: {batch_status.batch.request_counts}")

# Retrieve the status of all batches
all_batches = client.batches.list()

# Print all batch statuses
for batch in all_batches.data:
    print(f"Batch ID: {batch.id}")
    print(f"Status: {batch.status}")
    print(f"Created At: {batch.created_at}")
    print(f"Expires At: {batch.expires_at}")
    print(f"Request Counts: {batch.request_counts}")
    print("---")
```
Replace `"your_batch_id_here"` with the actual batch ID you want to check.

This code will connect to your local batch-gpt server and retrieve the status of either a specific batch or all batches. The response will include details such as the batch ID, status, creation time, expiration time, and request counts.

## Testing with Python Client

A Python test client is provided in the `test-python-client` directory.

1. Install the required Python package:
   ```
   cd test-python-client
   pip install -r requirements.txt
   ```

2. Run the test client:
   ```
   python client.py "Write a joke on Gandalf and Saruman"
   ```
   > **Note:** To effectively utilize batching, run multiple instances of the Python client simultaneously. This simulates concurrent requests, allowing the server to group them into batches for processing.

## Environment Variables

The following environment variables can be used to configure the application:

- `OPENAI_API_KEY`: Your OpenAI API key (required)
- `CLIENT_SERVING_MODE`: Set to "sync" for synchronous (default) or "async" for asynchronous serving mode
- `COLLATE_BATCHES_FOR_DURATION_IN_MS`: Duration to collate batches in milliseconds (default: 5000)
- `COLLECT_BATCH_STATS_POLLING_MAX_INTERVAL_SECONDS`: Maximum interval (in seconds) between polling attempts when collecting batch statistics. This value caps the exponential backoff for long-running batches. Default is 300 seconds (5 minutes) if not set.
- `MONGO_HOST`: MongoDB server hostname (default: "localhost")
- `MONGO_PORT`: MongoDB server port (default: "27017")
- `MONGO_USER`: MongoDB username (default: "admin")
- `MONGO_PASSWORD`: MongoDB password (default: "password")
- `MONGO_DATABASE`: MongoDB database name (default: "batchgpt")

## Advanced Settings

Fine-tune Batch-GPT's behavior with these advanced configuration options for optimal performance in various scenarios.

### Serving Modes

Batch-GPT supports two serving modes:

1. Synchronous Mode (Default):
   - Similar to the standard OpenAI requests, clients remain blocked after making a request to the server.
   - Ideal for low-volume scenarios where
   - Set `CLIENT_SERVING_MODE=sync` or leave unset

2. Asynchronous Mode:
   - Returns immediately with a submission confirmation
   - Ideal for high-volume scenarios where each worker remaining blocked on a response is not practical
   - Set `CLIENT_SERVING_MODE=async`

To change the serving mode, set the `CLIENT_SERVING_MODE` environment variable before starting the server.

### Batch Statistics Polling

The server uses an exponential backoff strategy when polling for batch statistics to reduce unnecessary API calls for long-running batches. The `COLLECT_BATCH_STATS_POLLING_MAX_INTERVAL_SECONDS` environment variable sets an upper limit on this interval.

For example:
```bash
export COLLECT_BATCH_STATS_POLLING_MAX_INTERVAL_SECONDS=600
```
This would set the maximum polling interval to 10 minutes. The actual polling interval starts smaller and increases exponentially up to this maximum value.

## Development

### Project Structure

- `/server`: Contains the main server code
  - `/db`: Database interactions
  - `/handlers`: HTTP request handlers
  - `/logger`: Custom logging setup
  - `/models`: Data models
  - `/services`: Business logic and OpenAI interactions
- `/local/mongo`: Docker setup for local MongoDB instance
- `/test-python-client`: Python client for testing the server

### Adding New Features

1. Implement new functionality in the appropriate package under `/server`
2. Update handlers in `/server/handlers` if adding new endpoints
3. Modify `/server/main.go` to wire up new endpoints or services
4. Update this README with any new setup or usage instructions

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
