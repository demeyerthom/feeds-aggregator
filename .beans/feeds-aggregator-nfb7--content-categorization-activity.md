---
# feeds-aggregator-nfb7
title: Content categorization activity
status: completed
type: feature
priority: normal
created_at: 2026-03-11T16:25:58Z
updated_at: 2026-03-12T20:53:10Z
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

## Summary of Work

- Task: Update FeedItemDocument type to include Categories field - Added Categories []string field with bson tag to FeedItemDocument struct
- Task: Create categorization prompt - Created internal/prompt/categorization_prompt.go with BuildCategorizationPrompt function and 10-category taxonomy as examples
- Task: Create CategorizeContent activity - Created internal/activity/categorize_content.go following closure pattern, parses JSON response for 1-5 categories, updates MongoDB
- Task: Register CategorizeContent activity in worker - Added activity registration in cmd/worker/main.go using same dependencies as CreateSummary
- Task: Update workflow to include categorization step - Added CategorizeContent activity call after CreateSummary in IngestFeedItem workflow
