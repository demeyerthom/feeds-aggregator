# Build stage
FROM golang:1.25-alpine AS builder

# Install CA certificates
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Build argument to specify which command to build (worker, ingester, processor)
ARG COMMAND

COPY go.mod go.sum ./

# Download dependencies with cache mount
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY ./internal ./internal
COPY "./cmd" "./cmd"

# Build with cache mount for faster rebuilds
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app cmd/${COMMAND}/main.go

# Final stage
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/app /app
ENTRYPOINT ["/app"]

