---
# feeds-aggregator-bbiz
title: Update docker-compose.yaml environment variables
status: completed
type: task
priority: normal
created_at: 2026-03-11T19:43:30Z
updated_at: 2026-03-11T19:56:58Z
parent: feeds-aggregator-65zw
---

Description: Update worker service environment in docker-compose.yaml to use new provider config structure. Remove old OLLAMA_HOST/OLLAMA_MODEL and ZEN_* variables, add new provider-specific variables.

Output Requirements:
- Worker service environment updated with new variable names
- Ollama provider configured with OLLAMA_HOST=http://ollama:11434 and OLLAMA_MODEL=qwen3.5:27b
- Remove stale ZEN_API_KEY, ZEN_MODEL variables (old OLLAMA_* vars already removed in current file)

Acceptance Criteria:
- docker-compose.yaml valid YAML
- Environment variables match new config struct field names
- OLLAMA_HOST set to http://ollama:11434
- OLLAMA_MODEL set to qwen3.5:27b

Context & Research:
- Current docker-compose.yaml lines 19-36
- Worker service environment section
- Keep OTEL_HOST, REDIS_HOST, TEMPORAL_HOST, LOG_LEVEL as-is
- Current file already has OLLAMA_HOST and OLLAMA_MODEL (lines 30-31) - update values, don't add new vars

Open Questions: None

Dependencies: None

## Summary of Changes\n\nUpdated docker-compose.yaml worker service environment:\n- Changed OLLAMA_MODEL from `${OLLAMA_MODEL:-gemma3}` to `qwen3.5:27b`\n- OLLAMA_HOST already correctly set to `http://ollama:11434`\n- No ZEN_* variables present (already removed)\n- YAML validated successfully

## Coder Notes\n- Updated OLLAMA_MODEL from gemma3 to qwen3.5:27b\n- OLLAMA_HOST already correct at http://ollama:11434\n- No ZEN_* variables present (already removed)\n- YAML validated successfully
