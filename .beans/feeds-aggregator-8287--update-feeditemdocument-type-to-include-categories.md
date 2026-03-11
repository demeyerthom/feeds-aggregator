---
# feeds-aggregator-8287
title: Update FeedItemDocument type to include Categories field
status: todo
type: task
created_at: 2026-03-11T16:26:07Z
updated_at: 2026-03-11T16:26:07Z
parent: feeds-aggregator-nfb7
---

Description: Add a Categories []string field to the FeedItemDocument struct in internal/types.go with the bson tag "categories". This will store the categories returned by the LLM.

Output Requirements: The FeedItemDocument struct in internal/types.go should have a Categories []string field with bson:"categories" tag.

Acceptance Criteria:
- FeedItemDocument struct has Categories field
- Field has correct bson tag
- Code compiles successfully
