# Ingester Service

The Ingester is a feed polling service that discovers new articles from configured RSS/Atom feeds and triggers processing workflows.

## Overview

**Service Name**: `feeds-ingester`  
**Entry Point**: `cmd/ingester/main.go`

## What the Ingester Does

1. **Polls Feeds**: Regularly fetches from configured RSS/Atom feeds
2. **Checks for Duplicates**: Uses Redis to avoid processing the same article twice
3. **Starts Processing**: Triggers a workflow in the Worker for each new article

## How It Works

```
Every minute (configurable):
         │
         ▼
┌─────────────────────┐
│ Fetch all feeds    │ ← Parse RSS/Atom
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│ For each article:  │
│   - Check if seen   │ ← Redis lookup
│   - If new:         │
│     - Mark seen     │ ← 7-day TTL in Redis
│     - Start workflow│ ← Tell Worker to process
└─────────────────────┘
```

The Ingester loads feed URLs from a JSON file. Each article link is checked against Redis—if it's new, a workflow is started and the link is stored with a 7-day expiration.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `REDIS_HOST` | `localhost:6379` | Redis server |
| `TEMPORAL_HOST` | `localhost:7233` | Temporal server |
| `SOURCE_FILE` | `/feeds.json` | Feed list JSON file |
| `TICKER_INTERVAL` | `1m` | How often to poll |

## Feed List

Feeds are defined in a JSON file:

```json
[
  { "title": "Hacker News", "xmlUrl": "https://news.ycombinator.com/rss" },
  { "title": "TechCrunch", "xmlUrl": "https://techcrunch.com/feed/" }
]
```
