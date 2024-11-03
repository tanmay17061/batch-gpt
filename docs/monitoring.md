---
layout: default
title: Monitoring
nav_order: 8
---

# Monitoring
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Batch Monitor Tool

### Starting the Monitor
```bash
./batch-monitor
```

### Interface Features

#### Navigation
- ↑/↓: Navigate batches
- ←/→: Switch tabs
- Tab: Cycle through tabs
- Enter: View details
- q: Quit

#### Status Tabs
- Active Batches
- Completed Batches
- Failed Batches
- Expired Batches

#### Display Elements
- Batch ID
- Status
- Creation Time
- Progress Bar
- Request Counts

## API Status Endpoints

### Single Batch Status
```bash
curl http://localhost:8080/v1/batches/{batch_id}
```

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

### All Batches Status
```bash
curl http://localhost:8080/v1/batches
```

## Logging

### Log Levels
- INFO: General operation info
- WARNING: Non-critical issues
- ERROR: Critical problems

### Log Location
- Standard output (stdout)
- Error output (stderr)
- MongoDB logs (for persistence)

## Performance Metrics

### Batch Statistics
- Total requests
- Completed requests
- Success rate
- Processing time

### Cache Performance
- Hit rate
- Miss rate
- Cache size
