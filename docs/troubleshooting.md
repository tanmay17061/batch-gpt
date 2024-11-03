---
layout: default
title: Troubleshooting
nav_order: 10
---

# Troubleshooting
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Common Issues

### Connection Problems

#### MongoDB Connection Failed
```
ERROR: Failed to connect to MongoDB
```

Solutions:
1. Check MongoDB is running:
   ```bash
   docker ps | grep mongodb
   ```
2. Verify connection settings:
   ```bash
   echo $MONGO_HOST
   echo $MONGO_PORT
   ```
3. Test connection manually:
   ```bash
   mongosh -u $MONGO_USER -p $MONGO_PASSWORD
   ```

#### Server Won't Start
```
ERROR: address already in use
```

Solutions:
1. Check for running instances:
   ```bash
   ps aux | grep batch-gpt
   ```
2. Kill existing process:
   ```bash
   kill -9 <PID>
   ```

### Processing Issues

#### Long Response Times
```
WARNING: Batch processing exceeding expected duration
```

Solutions:
1. Check OpenAI batch status
2. Verify network connectivity
3. Review batch window setting:
   ```bash
   echo $COLLATE_BATCHES_FOR_DURATION_IN_MS
   ```

#### Cache Problems
```
ERROR: Failed to cache response
```

Solutions:
1. Check MongoDB space:
   ```bash
   db.stats()
   ```
2. Verify cache collection:
   ```bash
   use batchgpt
   db.cached_responses.stats()
   ```

## Diagnostic Tools

### Log Analysis
```bash
# View recent logs
tail -f logs/batch-gpt.log

# Search for errors
grep ERROR logs/batch-gpt.log
```

### MongoDB Diagnostics
```bash
# Check collections
use batchgpt
show collections

# View batch status
db.batch_logs.find().sort({timestamp: -1}).limit(1)
```

### Monitor Tool
```bash
# Start monitor
./batch-monitor

# Filter for failed batches
# Use arrow keys to navigate to "Failed" tab
```

## Performance Optimization

### Batch Window Tuning
- Start with 5000ms
- Monitor batch sizes
- Adjust based on volume:
  ```bash
  export COLLATE_BATCHES_FOR_DURATION_IN_MS=3000  # More frequent, smaller batches
  ```

### Cache Management
- Regular maintenance
- Index optimization
- Storage monitoring

## Getting Help

1. Check documentation
2. Review GitHub issues
3. Open new issue with:
   - Error messages
   - Configuration
   - Logs
   - Steps to reproduce
