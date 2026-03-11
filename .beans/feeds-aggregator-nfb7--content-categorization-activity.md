---
# feeds-aggregator-nfb7
title: Content categorization activity
status: completed
type: feature
priority: normal
created_at: 2026-03-11T16:25:58Z
updated_at: 2026-03-11T17:41:30Z
parent: feeds-aggregator-vxnm
---

Add a new Temporal activity that sends article content to an LLM and receives 1-5 category tags. The categories are stored in the MongoDB document alongside the summary.

## General Requirements
- Define a category taxonomy (10 categories provided by user)
- Create a prompt instructing the LLM to return categories as JSON
- Parse LLM response and update MongoDB document
- Integrate into existing workflow after summary creation

## Design Choices
- Reuse existing LLM client (Zen/OpenAI) for consistency
- Follow existing closure-based activity pattern
- Categories stored as string array in MongoDB
