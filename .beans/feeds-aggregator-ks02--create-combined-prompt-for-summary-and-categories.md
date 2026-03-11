---
# feeds-aggregator-ks02
title: Create combined prompt for summary and categories
status: todo
type: task
created_at: 2026-03-11T19:50:02Z
updated_at: 2026-03-11T19:50:02Z
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
