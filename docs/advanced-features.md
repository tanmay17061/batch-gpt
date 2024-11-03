---
layout: default
title: Advanced Features
nav_order: 6
---

# Advanced Features
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Serving Modes

### Synchronous Mode
```bash
export CLIENT_SERVING_MODE=sync  # Default
```
- Blocks until response is available
- Similar to standard OpenAI API
- Best for low-volume scenarios

### Asynchronous Mode
```bash
export CLIENT_SERVING_MODE=async
```
- Returns immediately with submission confirmation
- Suitable for high-volume applications
- Requires separate status checking

### Cache-Only Mode
```bash
export CLIENT_SERVING_MODE=cache
```
- Serves only cached responses
- No new API calls
- Processes pending batches

## Caching System

### Cache Configuration
- Automatic request hashing
- MongoDB-based storage
- Cross-session persistence

### Cache Operations
```python
# Cache hit example
response1 = client.chat.completions.create(...)
response2 = client.chat.completions.create(...)  # Same request returns cached response
```

## Batch Recovery

### Automatic Recovery
- Detects interrupted batches
- Resumes processing on restart
- Updates original requesters

### Manual Recovery
```bash
# Check dangling batches
python client.py --api status_all_batches --status_filter not_completed
```

## Advanced Monitoring

### Custom Polling Intervals
```bash
export COLLECT_BATCH_STATS_POLLING_MAX_INTERVAL_SECONDS=600
```

### Progress Tracking
- Real-time completion statistics
- Request counts monitoring
- Error tracking

## Performance Tuning

### Batch Window Optimization
```bash
# Adjust batch collection window
export COLLATE_BATCHES_FOR_DURATION_IN_MS=3000  # 3 seconds
```

### MongoDB Optimization
- Index management
- Connection pooling
- Query optimization
