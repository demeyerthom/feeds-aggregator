---
# feeds-aggregator-j96n
title: Update CategorizeContent activity to use configurable text limit
status: completed
type: task
priority: normal
created_at: 2026-03-11T17:34:23Z
updated_at: 2026-03-11T19:12:01Z
parent: feeds-aggregator-vm05
blocked_by:
    - feeds-aggregator-qjl0
---

Description: Modify CategorizeContent activity to accept text limit config and pass it to the text extractor.

Output Requirements: Activity accepts text limit parameter and passes it to textextractor.ExtractArticleText

Acceptance Criteria: CategorizeContent uses configurable text limit from config

Context & Research:
- See internal/activity/categorize_content.go
- See cmd/worker/main.go lines 222-227 for activity registration pattern
- Activity uses closure pattern to receive dependencies

Dependencies: Task 2 (text extractor updated)

## Coder Notes\n- Added textLimit parameter to CategorizeContent function\n- Changed ExtractArticleText call to use closure pattern with textLimit
