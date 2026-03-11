---
# feeds-aggregator-65zw
title: Mutually Exclusive Provider Configuration
status: completed
type: feature
priority: normal
created_at: 2026-03-11T19:43:20Z
updated_at: 2026-03-11T20:08:32Z
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

## Summary of Changes

### Configuration
- Replaced Zen config struct with two optional provider structs (Ollama and OpenCode)
- Each provider has Host, Model, and Enabled fields (OpenCode also has APIKey)
- Environment variables: OLLAMA_*, OPENCODE_*

### Validation
- Mutual exclusivity validation ensures exactly one provider is enabled
- Clear error messages guide users to enable the correct provider

### Client Initialization
- OpenAI-compatible client initialized for both providers
- Ollama: uses OLLAMA_HOST + /v1/ with placeholder API key
- OpenCode: uses OPENCODE_HOST with required API key
- Correct model parameter passed to activities

### Docker Compose
- Worker configured to use Ollama provider
- OLLAMA_ENABLED=true, OLLAMA_HOST=http://ollama:11434, OLLAMA_MODEL=qwen3.5:27b

### Commits
1. feat: replace Zen config with Ollama and OpenCode provider structs
2. feat: add provider validation logic with mutual exclusivity
3. feat: initialize OpenAI client for both Ollama and OpenCode providers
4. feat: add OLLAMA_ENABLED to docker-compose for provider selection
