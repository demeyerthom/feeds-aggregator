---
# feeds-aggregator-3ee4
title: Register CategorizeContent activity in worker
status: completed
type: task
priority: normal
created_at: 2026-03-11T16:26:52Z
updated_at: 2026-03-11T16:52:25Z
parent: feeds-aggregator-nfb7
blocked_by:
    - feeds-aggregator-ycb8
---

Description: Register the new CategorizeContent activity in cmd/worker/main.go following the same pattern as other activities. Import the activity and register it with the worker using the same LLM client and model configuration as CreateSummary.

Output Requirements:
- Modified cmd/worker/main.go
- Activity properly registered with worker

Acceptance Criteria:
- Worker compiles successfully
- Activity is registered with correct name (using GetFunctionName)
- Uses same Zen/OpenAI client as CreateSummary

## Summary of Changes

- Added CategorizeContent activity registration in cmd/worker/main.go (lines 222-227)
- Uses same dependencies as CreateSummary: feedItemCollection, zenClient, cfg.Zen.Model, cfg.Storage.HTMLDir
- Uses internal.GetFunctionName(internalactivity.CategorizeContent) for activity name
- Worker compiles successfully
