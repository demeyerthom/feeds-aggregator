---
# feeds-aggregator-vxnm
title: Add content categorization to feed item processing
status: completed
type: epic
priority: normal
created_at: 2026-03-11T16:25:54Z
updated_at: 2026-03-11T17:41:34Z
---

Add a new step to the Temporal workflow that uses an LLM to categorize blog posts based on their content. Categories are stored in MongoDB for later filtering and organization.

## User Requirements
- Use LLM to determine 1-5 categories for each blog post
- Categories based on predefined taxonomy (10 categories)
- Return categories as JSON array
- Store categories in MongoDB
