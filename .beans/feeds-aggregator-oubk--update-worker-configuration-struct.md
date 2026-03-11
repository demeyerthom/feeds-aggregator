---
# feeds-aggregator-oubk
title: Update worker configuration struct
status: completed
type: task
priority: normal
created_at: 2026-03-11T19:43:28Z
updated_at: 2026-03-11T19:58:32Z
parent: feeds-aggregator-65zw
---

Description: Replace current Zen struct in cmd/worker/main.go with two optional structs: Ollama (host, model) and OpenCode (host, model, api_key). Remove the existing Zen config block entirely.

Output Requirements:
- Configuration struct has two optional provider blocks
- Ollama struct with OLLAMA_HOST (default: http://localhost:11434) and OLLAMA_MODEL (default: gemma3)
- OpenCode struct with OPENCODE_HOST (default: https://opencode.ai/zen/v1), OPENCODE_MODEL (default: big-pickle), and OPENCODE_API_KEY (no default)
- Old Zen struct removed

Acceptance Criteria:
- go build ./cmd/worker succeeds
- Config struct compiles with env tags for all new fields

Context & Research:
- Current config in cmd/worker/main.go lines 41-69
- Uses github.com/Netflix/go-env for env var binding
- Follow existing pattern: nested anonymous structs with env tags

Open Questions: None

Dependencies: None

## Summary of Changes

- Replaced the Zen struct with two new optional provider structs:
  - Ollama: OLLAMA_HOST (default: http://localhost:11434), OLLAMA_MODEL (default: gemma3)
  - OpenCode: OPENCODE_HOST (default: https://opencode.ai/zen/v1), OPENCODE_MODEL (default: big-pickle), OPENCODE_API_KEY (no default)
- Updated all references to cfg.Zen to use cfg.OpenCode
- Updated error message from ZEN_API_KEY to OPENCODE_API_KEY
- Updated log message to include host information
- Build verified successful

## Coder Notes\n- Removed Zen struct and added Ollama and OpenCode provider structs\n- Ollama: OLLAMA_HOST (default: http://localhost:11434), OLLAMA_MODEL (default: gemma3)\n- OpenCode: OPENCODE_HOST (default: https://opencode.ai/zen/v1), OPENCODE_MODEL (default: big-pickle), OPENCODE_API_KEY (no default)\n- Updated all code references from cfg.Zen to cfg.OpenCode\n- Build verified successful
