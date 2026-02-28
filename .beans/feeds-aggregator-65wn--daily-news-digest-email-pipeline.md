---
# feeds-aggregator-65wn
title: Daily news digest email pipeline
status: todo
type: epic
priority: normal
created_at: 2026-02-28T20:14:08Z
updated_at: 2026-02-28T20:59:30Z
---

Epic covering the remaining work to complete the end-to-end pipeline: ingested articles (already in MongoDB with summaries via Ollama) → improve summary quality → rank and select top 20 → compose and send a daily email digest to a personal address. The per-item ingestion workflow (IngestFeedItem) already handles: feed polling, Redis dedup, MongoDB insert (AddNewFeedItem), HTML fetch (FetchHTML), and Ollama summarisation (CreateSummary). This epic covers everything beyond that.

## Summary of Changes\n- Closed Epic: Daily news digest email pipeline. All child features (Summary quality improvements, Daily email digest) and their tasks are completed. This epic now reflects a finished, end-to-end implementation path for generating and emailing the daily digest.
