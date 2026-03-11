---
# feeds-aggregator-vm05
title: Configurable Text Limit with OTEL Metrics
status: completed
type: feature
priority: normal
created_at: 2026-03-11T17:33:08Z
updated_at: 2026-03-11T19:19:45Z
---

Description: Replace the hardcoded 4000 character limit in the text extractor with a configurable value (default 400000), add warning logs when limit is reached, and emit OTEL histogram metrics for body character counts.\n\nGeneral Requirements:\n- Add environment variable TEXT_LIMIT with default 400000 in worker config\n- Update text extractor to accept limit configuration via closure pattern\n- Add slog warning when text is truncated\n- Add Int64Histogram metric for body character counts with custom bucket boundaries\n\nDesign Choices:\n- Use closure pattern to pass configuration and create metric internally\n- Bucket boundaries: [5000, 10000, 20000, 50000, 100000, 200000, 400000, 600000, 800000, 1000000]\n- Metric name: feeds.text_extractor.body_char_count

## Summary of Changes\n\n- Added TEXT_LIMIT config to worker (default 400000)\n- Updated text extractor to use closure pattern with configurable limit\n- Added Int64Histogram metric for body character counts with custom bucket boundaries\n- Added slog warning when text is truncated\n- Updated CreateSummary and CategorizeContent activities to use text limit\n- Updated tests for new closure pattern
