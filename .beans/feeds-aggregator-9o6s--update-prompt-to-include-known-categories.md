---
# feeds-aggregator-9o6s
title: Update prompt to include known categories
status: todo
type: task
created_at: 2026-03-11T20:37:47Z
updated_at: 2026-03-11T20:37:47Z
parent: feeds-aggregator-4kxe
---

Description: Modify BuildProcessContentPrompt to accept knownCategories parameter. Include known categories as "preferred" list in prompt, keeping taxonomy exemplars as style examples. Instruct LLM to prefer known categories but allow new ones.

Output Requirements:
- Update BuildProcessContentPrompt signature: BuildProcessContentPrompt(title, url, text string, knownCategories []string) string
- Add known categories section to prompt after taxonomy exemplars
- Instruct LLM: prefer known categories if good match, but can create new ones
- Update tests to pass knownCategories parameter

Acceptance Criteria:
- go build ./internal/prompt succeeds
- Prompt includes taxonomy exemplars AND known categories
- Clear instruction to prefer known categories
- Tests pass with new signature

Context & Research:
- Current prompt in internal/prompt/process_content_prompt.go
- Keep existing taxonomy exemplars (lines 10-22)
- Add new section: "KNOWN CATEGORIES (prefer these if good match):"
- Update process_content_prompt_test.go

Open Questions: None

Dependencies: None
