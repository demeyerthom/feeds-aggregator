# AGENTS.md — feeds-aggregator

**IMPORTANT**: before you do anything else, run the `beans prime` command and heed its output.

## Project Overview

Go-based RSS/Atom feed aggregation system using Temporal workflows. Two services:
**ingester** (polls feeds, deduplicates via Redis, triggers workflows) and
**worker** (executes Temporal activities: store to MongoDB, fetch HTML, summarize via Ollama).

Module: `github.com/demeyerthom/feeds-aggregator` | Go 1.25 | No CI/CD pipelines configured.

## General instructions

- If you are unsure how to do something, use `gh_grep` to search code examples from GitHub.
- Use context7 search  
- Target Go 1.25+. Use generics for type-safe data structures.
- Write tests with `testing` + `testify` for every function.
- Use Go modules. Binaries go in `bin/`.
- Follow the folder structure: `cmd/{app}/main.go`, `internal/`, `config/{app}/`, `docker/{app}/Dockerfile`.
- Include doc comments on all functions. Minimum: `@param`, `@return`, `@throws`, `@author`.
- Sanitize inputs, parameterize DB queries, enforce security best practices.
- Target Go 1.25+. Use generics for type-safe data structures.
- Write tests with `testing` + `testify` for every function.
- Use Go modules. Binaries go in `bin/`.
- Follow the folder structure: `cmd/{app}/main.go`, `internal/`, `config/{app}/`, `docker/{app}/Dockerfile`.
- Include doc comments on all functions. Minimum: `@param`, `@return`, `@throws`, `@author`.
- Sanitize inputs, parameterize DB queries, enforce security best practices.

## Build and Run Commands

```bash
# Build binaries (output to bin/)
CGO_ENABLED=0 go build -o bin/ingester cmd/ingester/main.go
CGO_ENABLED=0 go build -o bin/worker cmd/worker/main.go

# Run tests (none exist yet — see Testing section)
go test ./...                          # all tests
go test ./internal/...                 # internal package only
go test ./internal/activity/...        # single subpackage
go test -run TestFunctionName ./...    # single test by name
go test -v -count=1 ./internal/...     # verbose, no cache

# Vet and format
go vet ./...
gofmt -w .

# Dependencies
go mod tidy
go mod download

# Docker (requires external infrastructure network)
docker compose --profile remote up -d --build

# Task runner (deployment only)
task remote:deploy                     # sync config + docker compose up
task remote:sync                       # sync config files to remote
task remote:clean                      # remove containers
```

## Testing

**No test files exist yet.** When adding tests:
- Use the `testing` package with `github.com/stretchr/testify` for assertions (already a dependency).
- Place test files next to source: `internal/activity/add_new_feed_item_test.go`.
- Name test functions `TestXxx(t *testing.T)`.
- Activities use closures — test the returned function by injecting mock dependencies.
- External services (Redis, MongoDB, Ollama, Temporal) must be mocked or stubbed in tests.

## Project Structure

```
feeds-aggregator/
├── cmd/
│   ├── ingester/main.go          # Feed polling service entry point
│   └── worker/main.go            # Temporal worker entry point
├── internal/
│   ├── activity/                  # Temporal activity implementations
│   │   ├── add_new_feed_item.go   # Insert feed item into MongoDB
│   │   ├── fetch_html.go          # Fetch and store HTML to disk
│   │   └── create_summary.go      # Ollama summarization and MongoDB update
│   ├── workflow/
│   │   └── ingest_feed_item.go    # Orchestrates the 3 activities in sequence
│   ├── types.go                   # Shared types: Feed, FeedItem, FeedItemDocument, FeedList
│   ├── constants.go               # Shared constants: TaskQueueName, Mongo DB/collection names
│   ├── logger.go                  # MultiHandler, ParseLogLevel
│   ├── otel.go                    # OpenTelemetry SDK bootstrap
│   ├── propagator.go              # Temporal context propagator
│   └── utils.go                   # GetFunctionName helper
├── config/ingester/feeds.json     # Feed source list
├── docker/Dockerfile              # Multi-stage build (golang:1.25-alpine to scratch)
├── docker-compose.yaml            # Services: ingester, worker
├── Taskfile.yaml                  # Remote deployment tasks only
├── data/                          # Fetched HTML files (gitignored)
└── bin/                           # Compiled binaries (gitignored)
```

## Code Style and Conventions

### Imports

Group imports in this order, separated by blank lines:
1. Standard library (`context`, `fmt`, `log/slog`, `os`, `time`, etc.)
2. Third-party packages (`github.com/...`, `go.temporal.io/...`, `go.opentelemetry.io/...`)

Internal package imports go with third-party. Use aliases when subpackage names collide:
```go
internalactivity "github.com/demeyerthom/feeds-aggregator/internal/activity"
internalworkflow "github.com/demeyerthom/feeds-aggregator/internal/workflow"
```

### Naming

- **Packages**: lowercase single words (`activity`, `workflow`, `internal`)
- **Exported functions**: PascalCase (`AddNewFeedItem`, `SetupOTelSDK`)
- **Constants**: PascalCase (`TaskQueueName`, `MongoDBName`)
- **Struct fields**: PascalCase with JSON/BSON/env tags
- **Config structs**: nested anonymous structs with `env:` tags using `github.com/Netflix/go-env`
- **Service names**: kebab-case strings (`"feeds-ingester"`, `"feeds-worker"`)

### Activity Pattern (Closure-Based Dependency Injection)

All Temporal activities use a closure pattern — the outer function accepts dependencies, returns the activity function:
```go
func ActivityName(dep1 Type1, dep2 Type2) func(ctx context.Context, input InputType) (OutputType, error) {
    return func(ctx context.Context, input InputType) (OutputType, error) {
        logger := activity.GetLogger(ctx)
        // implementation
    }
}
```

Activities are registered in `cmd/worker/main.go` with explicit names via `internal.GetFunctionName()`.

### Workflow Pattern

Workflows also use the closure pattern:
```go
func WorkflowName() func(ctx workflow.Context, input InputType) error {
    return func(ctx workflow.Context, input InputType) error {
        // orchestration logic
    }
}
```

### Error Handling

- Fatal errors during initialization: `slog.Error(msg, "err", err)` then `os.Exit(1)`
- Activity errors: log with `activity.GetLogger(ctx).Error(...)`, return the error
- Non-fatal processing errors: log and continue (e.g., feed parse failures return nil)
- Use `errors.Is()` for sentinel error checks, `errors.Join()` for combining errors
- Never panic — always return errors

### Logging

- Use `log/slog` (structured logging) exclusively — never `fmt.Println` or `log.Printf`
- Logger setup: `MultiHandler` combining JSON stdout + OpenTelemetry bridge
- In activities: use `activity.GetLogger(ctx)` (Temporal-aware)
- In workflows: use `workflow.GetLogger(ctx)` (Temporal-aware)
- Elsewhere: use `slog.Info/Error/Debug/Warn` with key-value pairs
- Format: `slog.Error("Description of what failed", "err", err, "key", value)`

### Configuration

- All config via environment variables using `github.com/Netflix/go-env`
- Config structs defined per-service in their `main.go` with `env:` tags
- Defaults provided in tags: `env:"VAR_NAME,default=value"`
- Config loaded in `init()` via `env.UnmarshalFromEnviron(&cfg)`

### Documentation

- Package-level doc comments on `activity` and `workflow` packages
- Function doc comments with `@param`, `@return`, `@throws`, `@author` tags
- Keep comments concise — describe what and why, not how

### Types

- All shared types live in `internal/types.go`
- All shared constants live in `internal/constants.go`
- MongoDB documents use `bson:` struct tags
- JSON types use `json:` struct tags
- Use `primitive.ObjectID` for MongoDB IDs

## Key Dependencies

| Dependency | Purpose |
|---|---|
| `go.temporal.io/sdk` | Workflow orchestration |
| `github.com/redis/go-redis/v9` | Feed item deduplication (24h x 7 TTL) |
| `go.mongodb.org/mongo-driver` | Feed item persistence |
| `github.com/mmcdole/gofeed` | RSS/Atom feed parsing |
| `github.com/ollama/ollama/api` | LLM summarization |
| `go.opentelemetry.io/otel` | Distributed tracing, metrics, logging |
| `github.com/Netflix/go-env` | Environment variable config binding |
| `github.com/stretchr/testify` | Test assertions (indirect, use for new tests) |

## Infrastructure Requirements

Services need: Redis, MongoDB, Temporal server, Ollama, OpenTelemetry Collector.
All services connect via Docker network `infrastructure` (external).
Docker builds use multi-stage: `golang:1.25-alpine` then `scratch` (no shell in prod image).
