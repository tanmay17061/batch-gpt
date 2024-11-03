---
layout: default
title: Getting Started
nav_order: 2
---

# Getting Started
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Prerequisites

Before installing Batch-GPT, ensure you have:

- Go 1.23.0 or later
- Docker and Docker Compose
- An OpenAI API key
- MongoDB (included in docker-compose setup)

## Installation Options

### Option 1: Pre-compiled Binaries (Recommended)

1. Download the latest release for your system from the [Releases page](https://github.com/tanmay17061/batch-gpt/releases)

2. Extract the archive:
   ```bash
   tar xzf batch-gpt_[version]_[os]_[arch].tar.gz
   ```

3. Set up MongoDB:
   ```bash
   cd local/mongo
   docker-compose up -d
   ```

### Option 2: Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/tanmay17061/batch-gpt.git
   cd batch-gpt
   ```

2. Set up MongoDB:
   ```bash
   cd local/mongo
   docker-compose up -d
   cd ../..
   ```

3. Build the server:
   ```bash
   go build -o batch-gpt server/main.go
   ```

## Configuration

Set required environment variables:

```bash
export OPENAI_API_KEY=your_openai_api_key
export COLLATE_BATCHES_FOR_DURATION_IN_MS=5000
export MONGO_HOST=localhost
export MONGO_PORT=27017
export MONGO_USER=admin
export MONGO_PASSWORD=password
export MONGO_DATABASE=batchgpt
```

## Running the Server

Start the server:
```bash
./batch-gpt
```

The server will start on `http://localhost:8080`.

## Verification

Test the installation using the provided Python client:

```bash
cd test-python-client
pip install -r requirements.txt
python client.py --api chat_completions --content "Hello, World!"
```

## Next Steps

- Learn about different [serving modes](advanced-features#serving-modes)
- Set up the [batch monitor](monitoring)
- Configure [caching](advanced-features#caching)
