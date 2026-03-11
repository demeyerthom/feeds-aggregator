---
# feeds-aggregator-65zw
title: Mutually Exclusive Provider Configuration
status: in-progress
type: feature
priority: normal
created_at: 2026-03-11T19:43:20Z
updated_at: 2026-03-11T19:54:34Z
---

Description: Add two optional configuration structs - Ollama (host, model) and OpenCode (host, model, api_key). Worker validates that exactly one is configured.

General Requirements:
- Two optional config blocks: Ollama and OpenCode
- Ollama struct: OLLAMA_HOST, OLLAMA_MODEL
- OpenCode struct: OPENCODE_HOST, OPENCODE_MODEL, OPENCODE_API_KEY
- Validation: exactly one provider must be configured (error if 0 or 2)
- Use OpenAI-compatible API for both (Ollama has /v1/ endpoint)

Design Choices:
- Mutually exclusive config blocks make the intent explicit
- Leverage Ollama's OpenAI-compatible API to use same OpenAI Go SDK for both
- Single client interface with different base URLs based on which provider is configured
