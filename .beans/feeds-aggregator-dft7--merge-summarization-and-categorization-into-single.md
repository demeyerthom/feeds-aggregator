---
# feeds-aggregator-dft7
title: Merge Summarization and Categorization into Single Activity
status: completed
type: feature
priority: normal
created_at: 2026-03-11T19:49:56Z
updated_at: 2026-03-11T20:27:25Z
---

Description: Combine the two LLM activities into one, reducing API calls and simplifying the workflow. The LLM will return a single JSON response containing both the summary and categories.

General Requirements:
- Single LLM call returns JSON with both summary and categories
- JSON format: {"summary": "...", "categories": ["cat1", "cat2"]}
- Single MongoDB update operation for both fields
- Reduce API calls from 2 to 1 per feed item

Design Choices:
- New ProcessContent activity replaces CreateSummary and CategorizeContent
- Combined prompt includes summarization instructions and categorization taxonomy
- Graceful JSON parsing error handling

## Summary of Changes

### Combined Prompt
- Created BuildProcessContentPrompt function in internal/prompt/process_content_prompt.go
- Combines summarization (2-3 sentences) and categorization (1-5 categories)
- JSON output format: {"summary": "...", "categories": ["..."]}
- Includes all 10 taxonomy exemplars

### ProcessContent Activity
- Created internal/activity/process_content.go
- Single LLM call for both summary and categories
- JSON response parsing with validation (1-5 categories)
- Single MongoDB update for both fields

### Workflow Update
- Updated internal/workflow/ingest_feed_item.go
- Replaced CreateSummary and CategorizeContent with single ProcessContent call
- Workflow now executes 3 activities instead of 4

### Worker Registration & Cleanup
- Updated cmd/worker/main.go to register ProcessContent
- Removed CreateSummary and CategorizeContent registrations
- Deleted 7 old files (activities and prompts)

### Benefits
- 50% reduction in LLM API calls (from 2 to 1 per feed item)
- Simpler workflow (3 activities instead of 4)
- Single MongoDB update for both summary and categories

### Commits
1. feat: create combined prompt for summary and categories
2. feat: create ProcessContent activity combining summary and categorization
3. feat: update workflow to use single ProcessContent activity
4. feat: update worker registration and remove old activities
