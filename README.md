# Feeds Aggregator

A distributed RSS/Atom feed aggregator and article processor built with Go, Temporal, and AI-powered summarization.

## Overview

Feeds Aggregator is a scalable system that automatically monitors RSS and Atom feeds, processes new articles, and generates AI-powered summaries. The system uses Temporal for reliable workflow orchestration, MongoDB for storage, Redis for deduplication, and Ollama for generating article summaries with large language models.

## Architecture

The system consists of two main components:

- **Ingester**: Periodically polls configured RSS/Atom feeds, detects new articles, and initiates processing workflows
- **Worker**: Executes Temporal workflows to fetch article HTML, store content, and generate AI summaries

## Features

- 🔄 Automatic feed polling and new article detection
- 📥 HTML content fetching and storage
- 🤖 AI-powered article summarization using Ollama (LLM)
- 🔁 Reliable workflow orchestration with Temporal
- 📊 Comprehensive observability with OpenTelemetry
- 💾 MongoDB storage for articles and metadata
- ⚡ Redis-based deduplication
- 🐳 Docker Compose deployment ready

## Technology Stack

- **Go 1.25+** - Core application language
- **Temporal** - Workflow orchestration
- **MongoDB** - Document storage
- **Redis** - Caching and deduplication
- **Ollama** - Local LLM for summarization
- **OpenTelemetry** - Observability and tracing

## Documentation

See the `/docs` folder for detailed documentation:

- [Architecture](docs/architecture.md) - System design and data flow
- [Ingester](docs/ingester.md) - Feed polling service
- [Worker](docs/worker.md) - Processing pipeline service