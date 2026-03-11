---
# feeds-aggregator-mzip
title: Update CreateSummary activity to use configurable text limit
status: completed
type: task
priority: normal
created_at: 2026-03-11T17:34:11Z
updated_at: 2026-03-11T19:11:49Z
parent: feeds-aggregator-vm05
blocked_by:
    - feeds-aggregator-qjl0
---

Description: Modify CreateSummary activity to accept text limit config and pass it to the text extractor.

Output Requirements: Activity accepts text limit parameter and passes it to textextractor.ExtractArticleText

Acceptance Criteria: CreateSummary uses configurable text limit from config

Context & Research:
- See internal/activity/create_summary.go
- See cmd/worker/main.go lines 216-221 for activity registration pattern
- Activity uses closure pattern to receive dependencies

Dependencies: Task 2 (text extractor updated)

## Coder Notes\n- Added textLimit parameter to CreateSummary function\n- Changed ExtractArticleText call to use closure pattern with textLimit
