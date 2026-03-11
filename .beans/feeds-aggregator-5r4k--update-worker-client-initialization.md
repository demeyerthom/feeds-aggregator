---
# feeds-aggregator-5r4k
title: Update worker client initialization
status: todo
type: task
created_at: 2026-03-11T19:43:57Z
updated_at: 2026-03-11T19:43:57Z
parent: feeds-aggregator-65zw
blocked_by:
    - feeds-aggregator-3pw8
---

Description: Initialize OpenAI client with correct base URL based on configured provider. For Ollama use OLLAMA_HOST + /v1/, for OpenCode use OPENCODE_HOST. API key required only for OpenCode.

Output Requirements:
- Single openai.Client variable initialized based on provider
- Ollama: baseURL = OLLAMA_HOST + /v1/, apiKey = ollama (placeholder, ignored)
- OpenCode: baseURL = OPENCODE_HOST, apiKey = OPENCODE_API_KEY (required)
- Model name passed to activities comes from configured provider

Acceptance Criteria:
- go build ./cmd/worker succeeds
- Worker initializes client correctly for both provider types
- Activities receive correct model name from config

Context & Research:
- Current client init at cmd/worker/main.go lines 175-184
- Uses github.com/openai/openai-go/v3 with option.WithAPIKey and option.WithBaseURL
- Ollama OpenAI compatibility: http://localhost:11434/v1/ (see context7 docs)
- Activities CreateSummary and CategorizeContent need model parameter

Open Questions: None

Dependencies: This task depends on feeds-aggregator-3pw8 being completed first
