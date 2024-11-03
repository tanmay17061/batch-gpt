---
layout: default
title: API Reference
nav_order: 9
---

# API Reference
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Chat Completions

### Create Chat Completion

`POST /v1/chat/completions`

Request:
```json
{
  "model": "gpt-3.5-turbo",
  "messages": [
    {
      "role": "user",
      "content": "Hello!"
    }
  ]
}
```

Response:
```json
{
  "id": "chatcmpl-123",
  "object": "chat.completion",
  "created": 1677652288,
  "choices": [{
    "index": 0,
    "message": {
      "role": "assistant",
      "content": "Hello! How can I help you today?"
    },
    "finish_reason": "stop"
  }],
  "usage": {
    "prompt_tokens": 9,
    "completion_tokens": 12,
    "total_tokens": 21
  }
}
```

## Batch Operations

### Retrieve Batch

`GET /v1/batches/{batch_id}`

Response:
```json
{
  "batch": {
    "id": "batch_123",
    "status": "completed",
    "created_at": 1678901234,
    "expires_at": 1678987634,
    "request_counts": {
      "total": 10,
      "completed": 10
    }
  }
}
```

### List Batches

`GET /v1/batches`

Response:
```json
{
  "data": [
    {
      "batch": {
        "id": "batch_123",
        "status": "completed",
        "created_at": 1678901234,
        "expires_at": 1678987634,
        "request_counts": {
          "total": 10,
          "completed": 10
        }
      }
    }
  ]
}
```

## Error Handling

### Error Response Format
```json
{
  "error": {
    "type": "invalid_request_error",
    "message": "Description of the error"
  }
}
```

### Common Error Types
- `invalid_request_error`
- `internal_server_error`
- `batch_processing_error`
