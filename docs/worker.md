# Worker Service

The Worker service executes the feed item processing pipeline. It receives workflows triggered by the Ingester and processes articles through to completion.

## Overview

**Service Name**: `feeds-worker`  
**Entry Point**: `cmd/worker/main.go`

## What the Worker Does

When the Ingester discovers a new feed item, it starts a workflow that the Worker executes in three steps:

1. **Insert to Database**: Creates a new document in MongoDB with the article link and title
2. **Fetch HTML**: Downloads the article's HTML content and saves it to disk
3. **Generate Summary**: Extracts text from the HTML, sends it to an LLM for summarization and categorization, then updates the MongoDB document with the results

## Processing Pipeline

```
New Feed Item (from Ingester)
         │
         ▼
┌─────────────────────┐
│ 1. Add to MongoDB   │ ← Store link and title
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│ 2. Fetch HTML      │ ← Download article content
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│ 3. Process Content │ ← LLM generates summary + categories
└─────────┬───────────┘
          │
          ▼
      Complete
```

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `MONGODB_URI` | `mongodb://localhost:27017` | MongoDB connection |
| `TEMPORAL_HOST` | `localhost:7233` | Temporal server |
| `HTML_STORAGE_DIR` | `./data` | Where to store fetched HTML |
| `TEXT_LIMIT` | `400000` | Max text characters to process |
| `OLLAMA_ENABLED` | `false` | Use Ollama for summaries |
| `OPENCODE_ENABLED` | `false` | Use OpenCode for summaries |

One LLM provider must be enabled (Ollama or OpenCode, not both).

## Data Storage

**MongoDB collection**: `feed_items`

| Field | Description |
|-------|-------------|
| `link` | Article URL (unique) |
| `title` | Article title |
| `summary` | AI-generated summary |
| `categories` | AI-generated categories |
| `created_at` | Insertion timestamp |
