---
# feeds-aggregator-ks02
title: Create combined prompt for summary and categories
status: completed
type: task
priority: normal
created_at: 2026-03-11T19:50:02Z
updated_at: 2026-03-11T20:11:50Z
parent: feeds-aggregator-dft7
---

Description: Create a new prompt function BuildProcessContentPrompt in internal/prompt/ that requests the LLM to return JSON with both summary and categories fields. Include the categorization taxonomy from the existing BuildCategorizationPrompt.

Output Requirements:
- New file internal/prompt/process_content_prompt.go
- Function BuildProcessContentPrompt(title, url, text string) string
- Prompt instructs LLM to return valid JSON: {"summary": "...", "categories": ["..."]}
- Include the taxonomy exemplars from categorization prompt
- Include summarization instructions from summary prompt

Acceptance Criteria:
- go build ./internal/prompt succeeds
- Prompt clearly instructs JSON-only output with no preamble
- Includes both summarization and categorization instructions

Context & Research:
- Existing prompts: internal/prompt/summary_prompt.go and internal/prompt/categorization_prompt.go
- Summary prompt: 2-3 sentence summary, no preamble
- Categorization prompt: 1-5 categories, taxonomy exemplars, JSON array output
- Combine both into single JSON object response

Open Questions: None

Dependencies: None

## Summary of Changes

Created new file internal/prompt/process_content_prompt.go with BuildProcessContentPrompt function that:
- Combines summarization and categorization instructions into a single prompt
- Instructs LLM to return JSON with both summary and categories fields
- Includes all 10 taxonomy exemplars from the categorization prompt
- Includes 2-3 sentence summary instructions from the summary prompt
- Explicitly instructs JSON-only output with no preamble

Also created test file internal/prompt/process_content_prompt_test.go with 5 test cases covering:
- Basic input inclusion
- Exemplar presence
- Summary instructions
- Categorization instructions
- JSON output format instructions

All tests pass.

## Coder Notes\n- Created internal/prompt/process_content_prompt.go with BuildProcessContentPrompt function\n- Created internal/prompt/process_content_prompt_test.go with 5 test cases\n- Prompt combines summarization (2-3 sentences) and categorization (1-5 categories)\n- JSON output format: {"summary": "...", "categories": ["..."]}\n- Includes all 10 taxonomy exemplars\n- Build and tests verified successful
