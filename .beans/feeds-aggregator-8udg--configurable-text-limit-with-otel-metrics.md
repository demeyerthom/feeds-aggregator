---
# feeds-aggregator-8udg
title: Configurable Text Limit with OTEL Metrics
status: scrapped
type: epic
priority: normal
created_at: 2026-03-11T17:33:00Z
updated_at: 2026-03-11T17:35:04Z
---

Make the text extractor limit configurable with OTEL metrics for monitoring text length distribution.\n\n- Add environment variable TEXT_LIMIT with default 400000\n- Update text extractor to accept configurable limit via closure pattern\n- Add slog warning when text is truncated\n- Emit Int64Histogram metric for body character counts with custom bucket boundaries\n\nBucket boundaries: [5000, 10000, 20000, 50000, 100000, 200000, 400000, 600000, 800000, 1000000]\nMetric name: feeds.text_extractor.body_char_count

## Reasons for Scrapping\n\nUser requested no epic - feature is sufficient standalone.
