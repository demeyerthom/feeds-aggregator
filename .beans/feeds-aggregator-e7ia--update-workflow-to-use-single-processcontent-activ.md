---
# feeds-aggregator-e7ia
title: Update workflow to use single ProcessContent activity
status: completed
type: task
priority: normal
created_at: 2026-03-11T19:50:20Z
updated_at: 2026-03-11T20:15:26Z
parent: feeds-aggregator-dft7
blocked_by:
    - feeds-aggregator-ak1k
---

Description: Update internal/workflow/ingest_feed_item.go to replace the two sequential activity calls (CreateSummary, CategorizeContent) with a single ProcessContent call.

Output Requirements:
- Remove CreateSummary activity execution (lines 47-52)
- Remove CategorizeContent activity execution (lines 54-59)
- Add single ProcessContent activity execution after FetchHTML
- Update workflow doc comment to reflect 3 activities instead of 4

Acceptance Criteria:
- go build ./internal/workflow succeeds
- Workflow executes 3 activities: AddNewFeedItem, FetchHTML, ProcessContent
- Workflow doc comment updated

Context & Research:
- Current workflow: internal/workflow/ingest_feed_item.go
- Activity execution pattern uses workflow.ExecuteActivity with internal.GetFunctionName
- Keep same ActivityOptions (5 minute timeout, 3 retry attempts)

Open Questions: None

Dependencies: This task depends on feeds-aggregator-ak1k being completed first

## Summary of Changes

Updated internal/workflow/ingest_feed_item.go:
- Updated doc comment from 'four activities' to 'three activities'
- Updated activity list in doc comment to: 'add feed item, fetch HTML, and process content'
- Removed CreateSummary activity execution
- Removed CategorizeContent activity execution
- Added single ProcessContent activity execution as the third activity
- Updated comment for third activity to 'process content (summary and categories)'

Build verification: go build ./internal/workflow succeeded

## Coder Notes\n- Updated doc comment from 'four activities' to 'three activities'\n- Replaced CreateSummary and CategorizeContent with single ProcessContent call\n- Updated activity list in doc comment\n- Build verified successful
