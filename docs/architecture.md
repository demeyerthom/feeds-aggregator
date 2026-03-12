# Architecture Overview

The Feeds Aggregator is a distributed system for automatically collecting, processing, and summarizing RSS/Atom feeds. It combines modern observability practices with reliable workflow orchestration to create a scalable feed processing pipeline.

## System Design

The system follows a **producer-consumer pattern** with two main services:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              FEEDS AGGREGATOR                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌──────────────────────┐                              ┌──────────────────┐ │
│  │       INGESTER       │                              │      WORKER      │ │
│  │                      │     Start Workflows      │                   │ │
│  │  • Polls RSS/Atom    │ ─────────────────────────▶  │  • Executes      │ │
│  │    feeds             │                              │    Activities    │ │
│  │  • Deduplicates      │                              │  • Manages       │ │
│  │    via Redis         │                              │    Workflows     │ │
│  │  • Triggers          │                              │                   │ │
│  │    Temporal          │                              │                   │ │
│  │    workflows         │                              │                   │ │
│  └──────────────────────┘                              └──────────────────┘ │
│           │                                                    │            │
│           │                                                    │            │
│           ▼                                                    ▼            │
│  ┌──────────────────────┐                              ┌──────────────────┐ │
│  │       REDIS         │                              │     MONGODB      │ │
│  │                      │                              │                   │ │
│  │  • Deduplication    │                              │  • Feed items    │ │
│  │    (7-day TTL)       │                              │  • Summaries     │ │
│  │                      │                              │  • Categories    │ │
│  └──────────────────────┘                              └──────────────────┘ │
│                                                                    │       │
│                                                                    │       │
│                                                          ┌─────────▼─────┐ │
│                                                          │   OLLAMA /    │ │
│                                                          │   OPENCODE    │ │
│                                                          │   (LLM)       │ │
│                                                          └───────────────┘ │
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                        TEMPORAL SERVER                               │  │
│  │                                                                       │  │
│  │  • Workflow state management                                        │  │
│  │  • Activity scheduling                                              │  │
│  │  • Reliable execution with retries                                  │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Data Flow

### 1. Feed Polling (Ingester)

The Ingester service periodically polls configured RSS/Atom feeds:

1. Loads feed URLs from `config/ingester/feeds.json`
2. Parses each feed using the `gofeed` library
3. For each feed item:
   - Checks Redis for duplicate detection (7-day TTL)
   - If new: stores the link in Redis and starts a Temporal workflow

### 2. Workflow Execution (Worker)

The Worker service executes the `IngestFeedItem` workflow which orchestrates three sequential activities:

1. **AddNewFeedItem**: Inserts the feed item into MongoDB, returns document with generated ID
2. **FetchHTML**: Fetches the article's HTML and saves it to disk
3. **ProcessContent**: Extracts text, sends to LLM for summarization and categorization, updates MongoDB

### 3. Observability

All services integrate with OpenTelemetry:

- **Traces**: Distributed tracing across services and activities
- **Metrics**: Custom metrics like `feeds.new_links` counter
- **Logs**: Structured JSON logging with OTel bridge

## Key Design Decisions

### Why Temporal?

Temporal provides:
- **Reliability**: Automatic retries with configurable policies
- **Visibility**: Workflow state and history for debugging
- **Scalability**: Decoupled activity execution across workers
- **Observability**: Built-in tracing integration

### Why Redis for Deduplication?

- Fast O(1) lookups for link checking
- TTL support (7 days) for automatic cleanup
- Low latency compared to MongoDB for high-frequency checks

### Why Two Services?

- **Ingester**: Lightweight polling, can run frequently without heavy resource use
- **Worker**: Can be scaled independently based on processing load
- Separation of concerns: discovery vs. processing

### Closure Pattern for Activities

Activities use dependency injection via closures:
```go
func Activity(dep1 Type1, dep2 Type2) func(ctx context.Context, input Input) (Output, error)
```

This allows:
- Easy mocking in tests
- Centralized dependency management in `main.go`
- No global state in activity packages

## Configuration

All configuration is environment-based:

| Variable | Service | Description |
|----------|---------|-------------|
| `REDIS_HOST` | Both | Redis connection string |
| `MONGODB_URI` | Worker | MongoDB connection string |
| `TEMPORAL_HOST` | Both | Temporal server address |
| `OTEL_HOST` | Both | OpenTelemetry collector |
| `OLLAMA_ENABLED` | Worker | Enable Ollama provider |
| `OPENCODE_ENABLED` | Worker | Enable OpenCode provider |
| `TICKER_INTERVAL` | Ingester | Feed polling interval |

## Infrastructure Requirements

- **Redis**: Feed item deduplication
- **MongoDB**: Persistent storage
- **Temporal**: Workflow orchestration
- **Ollama/OpenCode**: LLM for summarization
- **OpenTelemetry Collector**: Observability backend

All services run on the Docker network `infrastructure` (external).
