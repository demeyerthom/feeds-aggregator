---
# feeds-aggregator-dft7
title: Merge Summarization and Categorization into Single Activity
status: todo
type: feature
created_at: 2026-03-11T19:49:56Z
updated_at: 2026-03-11T19:49:56Z
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
