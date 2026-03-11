---
# feeds-aggregator-ycb8
title: Create CategorizeContent activity
status: todo
type: task
created_at: 2026-03-11T16:26:32Z
updated_at: 2026-03-11T16:26:32Z
parent: feeds-aggregator-nfb7
blocked_by:
    - feeds-aggregator-8287
    - feeds-aggregator-so6t
---

Description: Create a new activity in internal/activity/categorize_content.go that:
- Takes the FeedItemDocument (with ID and link)
- Reads the HTML file from disk (same file as used by CreateSummary)
- Extracts article text using existing text extractor
- Calls the LLM with the categorization prompt
- Parses the JSON response to extract categories (1-5)
- Updates the MongoDB document with the categories array

Follow the existing closure-based pattern used by CreateSummary. Create corresponding test file.

Output Requirements:
- New file internal/activity/categorize_content.go with CategorizeContent function
- New file internal/activity/categorize_content_test.go with tests
- Activity returns error if category parsing fails

Acceptance Criteria:
- Activity follows existing pattern (closure with dependencies)
- Correctly calls LLM and parses JSON response
- Updates MongoDB document with categories
- Tests pass
