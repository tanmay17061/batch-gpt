# batch-gpt
A jump-server to convert openai chat completion api requests to batched chat completion requests

## Features

- Accepts standard OpenAI API-compatible chat completion requests
- Batches multiple requests together before sending to OpenAI
- Provides a status endpoint for checking batch processing status
- Uses MongoDB for logging batch statuses

## Prerequisites

- Go 1.23.0 or later
- Docker and Docker Compose (for running MongoDB)
- OpenAI API key

## Setup

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
   ```

4. Build and run the server:
   ```
   go build -o batch-gpt server/main.go
   ./batch-gpt
   ```

The server will start on `http://localhost:8080`.

## Usage

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
    api_key="dummy_api_key",  # The API key is not used by batch-gpt, but is required by the client
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
    api_key="dummy_api_key",  # The API key is not used by batch-gpt, but is required by the client
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
