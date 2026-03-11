---
# feeds-aggregator-t5v8
title: Update worker registration and remove old activities
status: completed
type: task
priority: normal
created_at: 2026-03-11T19:50:26Z
updated_at: 2026-03-11T20:17:23Z
parent: feeds-aggregator-dft7
blocked_by:
    - feeds-aggregator-e7ia
---

Description: Update cmd/worker/main.go to register ProcessContent activity and remove CreateSummary and CategorizeContent registrations. Delete the old activity files and prompt files.

Output Requirements:
- Register ProcessContent activity in worker (w.RegisterActivityWithOptions)
- Remove CreateSummary registration
- Remove CategorizeContent registration
- Delete internal/activity/create_summary.go
- Delete internal/activity/categorize_content.go
- Delete internal/prompt/summary_prompt.go
- Delete internal/prompt/summary_prompt_test.go
- Delete internal/prompt/categorization_prompt.go
- Delete internal/prompt/categorization_prompt_test.go

Acceptance Criteria:
- go build ./cmd/worker succeeds
- Worker registers only: AddNewFeedItem, FetchHTML, ProcessContent
- Old files deleted
- go test ./... passes (after removing old test files)

Context & Research:
- Worker registration: cmd/worker/main.go lines 213-230
- Registration pattern: w.RegisterActivityWithOptions with internal.GetFunctionName
- Pass same dependencies to ProcessContent: feedItemCollection, client, model, dataDir, textLimit

Open Questions: None

Dependencies: This task depends on feeds-aggregator-e7ia being completed first

## Summary of Changes

- Updated cmd/worker/main.go to register ProcessContent activity instead of CreateSummary and CategorizeContent
- Deleted internal/activity/create_summary.go
- Deleted internal/activity/categorize_content.go
- Deleted internal/activity/categorize_content_test.go
- Deleted internal/prompt/summary_prompt.go
- Deleted internal/prompt/summary_prompt_test.go
- Deleted internal/prompt/categorization_prompt.go
- Deleted internal/prompt/categorization_prompt_test.go
- Added ErrInvalidCategoryCount error definition to process_content.go (was previously in deleted categorize_content.go)
- Build and tests pass successfully

## Coder Notes\n- Updated worker registration to use ProcessContent instead of CreateSummary and CategorizeContent\n- Deleted old activity files: create_summary.go, categorize_content.go, categorize_content_test.go\n- Deleted old prompt files: summary_prompt.go, summary_prompt_test.go, categorization_prompt.go, categorization_prompt_test.go\n- Moved ErrInvalidCategoryCount to process_content.go (was in deleted file)\n- Build and tests verified successful
