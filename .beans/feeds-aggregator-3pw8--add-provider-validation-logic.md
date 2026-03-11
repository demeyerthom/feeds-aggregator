---
# feeds-aggregator-3pw8
title: Add provider validation logic
status: todo
type: task
created_at: 2026-03-11T19:43:35Z
updated_at: 2026-03-11T19:43:35Z
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
