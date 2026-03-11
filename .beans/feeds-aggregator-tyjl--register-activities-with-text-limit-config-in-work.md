---
# feeds-aggregator-tyjl
title: Register activities with text limit config in worker
status: completed
type: task
priority: normal
created_at: 2026-03-11T17:34:29Z
updated_at: 2026-03-11T19:13:49Z
parent: feeds-aggregator-vm05
blocked_by:
    - feeds-aggregator-mzip
    - feeds-aggregator-j96n
---

Description: Update worker main.go to pass text limit config to the activities.

Output Requirements:
- Pass TEXT_LIMIT from config to CreateSummary activity
- Pass TEXT_LIMIT from config to CategorizeContent activity

Acceptance Criteria: Activities receive the text limit config correctly

Context & Research:
- See cmd/worker/main.go lines 209-227 for activity registration pattern
- Activities need to be updated to accept the text limit parameter

Dependencies: Tasks 3 and 4

## Coder Notes\n- Updated CreateSummary registration to pass cfg.TextExtractor.Limit\n- Updated CategorizeContent registration to pass cfg.TextExtractor.Limit
