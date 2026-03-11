---
# feeds-aggregator-e7ia
title: Update workflow to use single ProcessContent activity
status: todo
type: task
created_at: 2026-03-11T19:50:20Z
updated_at: 2026-03-11T19:50:20Z
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
