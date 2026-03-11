---
# feeds-aggregator-3pw8
title: Add provider validation logic
status: completed
type: task
priority: normal
created_at: 2026-03-11T19:43:35Z
updated_at: 2026-03-11T20:04:47Z
parent: feeds-aggregator-65zw
blocked_by:
    - feeds-aggregator-oubk
---

Description: After loading config in worker main.go, validate that exactly one provider is configured. Return clear error messages for invalid states.

Output Requirements:
- Validation function or inline check after env.UnmarshalFromEnviron
- Check if both Ollama and OpenCode are configured (error)
- Check if neither is configured (error)
- Clear error messages using slog.Error then os.Exit(1)

Acceptance Criteria:
- Worker fails to start with clear message if no provider configured
- Worker fails to start with clear message if both providers configured
- Worker starts successfully when exactly one provider is configured

Context & Research:
- Follow existing error handling pattern: slog.Error(msg, "err", err) then os.Exit(1)
- See cmd/worker/main.go lines 176-179 for existing validation pattern
- A provider is considered configured if its Model field is non-empty (or use a dedicated Enabled flag)

Open Questions: None

Dependencies: This task depends on feeds-aggregator-oubk being completed first



## Summary of Changes

Added provider validation logic to cmd/worker/main.go:

1. Added `Enabled` field to both Ollama and OpenCode config blocks with default=false
2. Added validation after MongoDB setup to check mutual exclusivity:
   - Error if both providers are enabled
   - Error if neither provider is enabled
   - Clear error messages using slog.Error then os.Exit(1)
3. Updated provider initialization to only initialize OpenCode client when OpenCode is enabled
4. Added log message when Ollama provider is selected

The validation ensures exactly one provider must be enabled via OLLAMA_ENABLED=true or OPENCODE_ENABLED=true environment variables.

## Coder Notes\n- Added Enabled fields to both provider config blocks (OLLAMA_ENABLED, OPENCODE_ENABLED, both default false)\n- Added mutual exclusivity validation with clear error messages\n- Updated provider initialization to be conditional based on enabled provider\n- OpenCode API key validation only runs when OpenCode is enabled\n- Build verified successful

## Review Findings\n- Reviewer found that zenClient is not initialized for Ollama provider\n- Reviewer found that wrong model parameter is passed for Ollama\n- These issues will be addressed by task 5r4k (Update worker client initialization)\n- Task 3pw8 scope was validation logic only; client initialization is task 5r4k's responsibility
