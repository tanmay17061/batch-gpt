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

Send POST requests to `/v1/chat/completions` with the same format as the OpenAI API. For example:

```bash
curl http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

### Checking Batch Status

To check the status of a batch, send a GET request to `/v1/batches/{batch_id}`:

```bash
curl http://localhost:8080/v1/batches/your_batch_id_here
```

## Testing with Python Client

A Python test client is provided in the `test-python-client` directory.

1. Install the required Python package:
   ```
   cd test-python-client
   pip install -r requirements.txt
   ```

2. Run the test client:
   ```
   python client.py "Your message here"
   ```

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
