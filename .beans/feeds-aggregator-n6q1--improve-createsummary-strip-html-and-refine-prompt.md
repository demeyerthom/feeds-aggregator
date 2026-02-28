---
# feeds-aggregator-n6q1
title: 'Improve CreateSummary: strip HTML and refine prompt before sending to Ollama'
status: completed
type: task
priority: normal
created_at: 2026-02-28T20:14:26Z
updated_at: 2026-02-28T20:58:59Z
parent: feeds-aggregator-9wbo
---

## Problem

The existing `CreateSummary` activity (`internal/activity/create_summary.go`) reads the raw HTML file from disk and sends **the entire HTML** — including `<script>`, `<style>`, navigation, footers, ads, and boilerplate — directly to Ollama. This wastes context window, degrades summary quality, and can exceed the model's token limit on large pages.

## What to do

1. **Add an HTML-to-text extraction step** before the Ollama call. Use `golang.org/x/net/html` (already transitively available) or a purpose-built library like `github.com/go-shiori/go-readability` to extract the article body text from the raw HTML. Strip all `<script>`, `<style>`, `<nav>`, `<header>`, `<footer>` elements and non-content boilerplate.
2. **Improve the Ollama prompt.** The current prompt is just `"Summarize this in 100 words or less:\n\n" + rawHTML`. Replace it with a system-prompt–style instruction that tells the model it is summarising a news article and should focus on the key facts, written in 2–3 sentences.
3. **Increase the activity timeout.** The workflow currently sets `StartToCloseTimeout: 30s` for all three activities. Ollama inference (especially on larger articles with a local model) can easily exceed 30 seconds. Give the `CreateSummary` activity its own `ActivityOptions` with a longer timeout (e.g. 2–5 minutes).
4. **Handle extraction failures gracefully.** If the readability extraction returns empty or fails, fall back to a naive tag-strip (regex or `html.Tokenizer`) rather than sending raw HTML.
5. **Add a max-length guard.** Truncate the extracted text to a sensible limit (e.g. 4 000 characters) before sending to Ollama so it fits comfortably within the model's context window.

## Existing code to modify

- `internal/activity/create_summary.go` — the activity itself
- `internal/workflow/ingest_feed_item.go` — give `CreateSummary` its own, longer activity options
- `go.mod` — add readability library dependency if chosen

## Acceptance criteria

- [ ] Raw HTML is converted to article-body plain text before being sent to Ollama
- [ ] Prompt is improved with clear instructions for news-article summarisation
- [ ] `CreateSummary` activity has a separate, longer timeout (≥ 2 min)
- [ ] Extraction failure falls back gracefully (no raw HTML sent to Ollama)
- [ ] Extracted text is truncated to a safe max length
- [ ] Unit tests cover: extraction logic, prompt construction, fallback path, truncation

## Summary of Changes\n- Implemented and closed the Summary Quality Improvements task as part of the epic. HTML extraction, trimmed text, improved prompt, and tests added.
