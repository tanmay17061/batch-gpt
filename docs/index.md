---
layout: default
title: Home
nav_order: 1
description: "Batch-GPT documentation home page"
permalink: /
---

# Batch-GPT Documentation
{: .fs-9 }

A jump-server to convert OpenAI chat completion API requests to batched requests
{: .fs-6 .fw-300 }

[Get Started](getting-started){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
[View on GitHub](https://github.com/tanmay17061/batch-gpt){: .btn .fs-5 .mb-4 .mb-md-0 }

---

## Overview

Batch-GPT enables seamless integration with OpenAI's Batch API by acting as a drop-in replacement for standard OpenAI API endpoints. It intelligently collects and batches requests, providing significant cost savings while maintaining compatibility with existing OpenAI clients.

### Quick Integration

It really is as simple as:
```diff
from openai import OpenAI
- client = OpenAI(api_key="sk-...")
+ client = OpenAI(api_key="dummy_openai_api_key", base_url="http://batch-gpt")
```
{: .highlight }

## Highlights

- **Cost Savings**: Up to 50% reduction using OpenAI's Batch API
- **Zero Config**: Drop-in replacement for OpenAI API clients
- **Reliability**: Automatic recovery of interrupted batches
- **Monitoring**: Real-time terminal-based batch status tracking
- **Flexibility**: Sync/Async/Cache operation modes

[Read More About Features](features){: .btn .btn-blue }
