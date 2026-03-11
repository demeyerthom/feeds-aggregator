---
# feeds-aggregator-t5v8
title: Update worker registration and remove old activities
status: todo
type: task
created_at: 2026-03-11T19:50:26Z
updated_at: 2026-03-11T19:50:26Z
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
