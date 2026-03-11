---
# feeds-aggregator-r6l3
title: Update workflow to include categorization step
status: completed
type: task
priority: normal
created_at: 2026-03-11T16:26:37Z
updated_at: 2026-03-11T16:52:25Z
parent: feeds-aggregator-nfb7
blocked_by:
    - feeds-aggregator-ycb8
---

Description: Modify the IngestFeedItem workflow in internal/workflow/ingest_feed_item.go to add a fourth activity call for CategorizeContent after the CreateSummary activity. Pass the feedItemDoc so the activity can access the MongoDB document ID and link.

Output Requirements:
- Modified internal/workflow/ingest_feed_item.go
- New activity call after CreateSummary
- Workflow still compiles and runs

Acceptance Criteria:
- Workflow compiles successfully
- CategorizeContent is called after CreateSummary
- Error handling is appropriate
