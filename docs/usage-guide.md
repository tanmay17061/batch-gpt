---
layout: default
title: Usage Guide
nav_order: 5
---

# Usage Guide
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Basic Usage

### Python Client

```python
from openai import OpenAI

client = OpenAI(
    api_key="dummy_openai_api_key",  # Any string works as the key
    base_url="http://localhost:8080/v1"
)

response = client.chat.completions.create(
    model="gpt-3.5-turbo",
    messages=[
        {"role": "user", "content": "Hello!"}
    ]
)

print(response.choices[0].message.content)
```

### curl Commands

Chat completion request:
```bash
curl http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

Check batch status:
```bash
curl http://localhost:8080/v1/batches/batch_123
```

List all batches:
```bash
curl http://localhost:8080/v1/batches
```

## Test Client Usage

The included Python test client provides easy testing:

```bash
# Send a chat completion request
python client.py --api chat_completions --content "Write a joke"

# Check specific batch status
python client.py --api status_single_batch --batch_id batch_123

# List all batches
python client.py --api status_all_batches

# List only completed batches
python client.py --api status_all_batches --status_filter completed
```

## Response Handling

### Synchronous Mode
```python
# Response will be returned when ready
response = client.chat.completions.create(...)
print(response.choices[0].message.content)
```

### Asynchronous Mode
```python
# Returns immediately with batch ID
response = client.chat.completions.create(...)
batch_id = response.id

# Check status later
status = client.batches.retrieve(batch_id)
print(f"Status: {status.batch.status}")
```

## Best Practices

1. **Request Batching**
   - Group similar requests together
   - Use appropriate batch window size
   - Consider request volume

2. **Error Handling**
   ```python
   try:
       response = client.chat.completions.create(...)
   except Exception as e:
       print(f"Error: {e}")
   ```

3. **Monitoring**
   - Use the batch monitor tool
   - Track batch statuses
   - Monitor cache hits/misses
