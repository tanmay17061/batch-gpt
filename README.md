# batch-gpt
A jump-server to convert openai chat completion api requests to batched chat completion requests

It really is as simple as
```diff
- client = OpenAI(api_key="sk-...")
+ client = OpenAI(api_key="dummy_openai_api_key", base_url="http://batch-gpt/v1")
```

## Features 🤓👍

- Ease of use: Accepts standard OpenAI-compatible APIs for chat completions, batch status, etc.
- Out-of-box cost reductions:
  1. ⭐ Uses OpenAI Batch API: Batch-GPT automatically converts chat completion requests to batch requests. OpenAI offers cost savings of up to 50% on their batch API.
  2. Automatic request caching: As of 10-10-2024, [OpenAI does not offer prompt caching on their batch API](https://platform.openai.com/docs/guides/prompt-caching/frequently-asked-questions). Batch-GPT caches all requests, and incurs zero costs for repeating requests.
- Increased reliability: In the scenario of a Batch-GPT server or client failure, all dangling-batches resume processing on the next server startup. Thanks to request caching, the client can make the same call to retrieve response from the cache.
- Above features are supported across service runs. Batch-GPT persists data on a MongoDB instance.
- Ability to view the completion status of all batches at once.
- Centralized OpenAI key distribution: No need to circulate your OpenAI key amongst users. Batch-GPT accepts `OPENAI_API_KEY`. Any client communicating with Batch-GPT will automatically assume this key.

## Limitations 🤔

- As of 10-10-2024, [OpenAI's Batch API claims an SLA of 24-hour turnaround time](https://platform.openai.com/docs/guides/batch/batch-api).
- While the reliability features listed above help mitigate effects of this issue, this high TAT might lead to delays in project/client-server disconnections... It is definitely not recommended for serving live requests!

> **💡** Explore OpenAI's [Realtime API](https://platform.openai.com/docs/guides/realtime) for such usecases.

## Prerequisites

- Go 1.23.0 or later
- Docker and (Docker Compose if running MongoDB through `local/mongo/docker-compose.yaml`)
- An OpenAI API key

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
