---
# feeds-aggregator-ak1k
title: Create ProcessContent activity
status: completed
type: task
priority: normal
created_at: 2026-03-11T19:50:10Z
updated_at: 2026-03-11T20:13:52Z
parent: feeds-aggregator-dft7
blocked_by:
    - feeds-aggregator-ks02
---

Description: Create a new activity internal/activity/process_content.go that combines summarization and categorization. Read HTML file, extract text, call LLM with combined prompt, parse JSON response, update MongoDB with both summary and categories in a single operation.

Output Requirements:
- New file internal/activity/process_content.go
- Function ProcessContent(c *mongo.Collection, client openai.Client, model, dataDir string, textLimit int) func(ctx context.Context, feedItemDoc internal.FeedItemDocument) error
- Use text extractor with textLimit (same pattern as existing activities)
- Call LLM with BuildProcessContentPrompt
- Parse JSON response into struct with Summary and Categories fields
- Single MongoDB UpdateOne with $set for both fields

Acceptance Criteria:
- go build ./internal/activity succeeds
- Activity follows closure pattern like existing activities
- Handles JSON parsing errors with clear logging
- Updates MongoDB document with both summary and categories

Context & Research:
- Existing activities pattern: internal/activity/create_summary.go and internal/activity/categorize_content.go
- Use textextractor.ExtractArticleText and textextractor.StripHTMLToPlainText
- MongoDB update: bson.M{"$set": bson.M{"summary": summary, "categories": categories}}
- Use encoding/json for parsing LLM response

Open Questions: None

Dependencies: This task depends on the combined prompt being created first

## Summary of Changes

Created internal/activity/process_content.go with the ProcessContent activity that:
- Follows the closure pattern like existing activities
- Reads HTML file from dataDir using feedItemDoc.ID.Hex()
- Extracts article text using textextractor.ExtractArticleText with fallback to StripHTMLToPlainText
- Calls LLM with prompt.BuildProcessContentPrompt for combined summarization and categorization
- Parses JSON response into processContentResponse struct with Summary and Categories fields
- Validates categories count (1-5)
- Updates MongoDB document with both summary and categories in a single UpdateOne operation

All acceptance criteria met:
- go build ./internal/activity succeeds
- Activity follows closure pattern
- Handles JSON parsing errors with clear logging
- Updates MongoDB document with both summary and categories

## Coder Notes\n- Created internal/activity/process_content.go with ProcessContent activity\n- Follows closure pattern like existing activities\n- Uses BuildProcessContentPrompt for combined LLM call\n- Parses JSON response into struct with Summary and Categories fields\n- Validates 1-5 categories (reuses ErrInvalidCategoryCount)\n- Single MongoDB UpdateOne for both fields\n- Build verified successful
