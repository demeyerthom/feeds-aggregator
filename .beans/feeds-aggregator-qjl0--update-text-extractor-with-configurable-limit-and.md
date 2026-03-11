---
# feeds-aggregator-qjl0
title: Update text extractor with configurable limit and metrics
status: todo
type: task
priority: normal
created_at: 2026-03-11T17:33:24Z
updated_at: 2026-03-11T17:34:32Z
parent: feeds-aggregator-vm05
blocked_by:
    - feeds-aggregator-e8t3
---

Description: Modify ExtractArticleText function to accept a limit parameter via closure pattern, log warning when truncated, and emit histogram metric for body character counts.

Output Requirements:
- Function signature changes to return a function that accepts limit
- Add slog.Warn when text is truncated (include original length and limit)
- Create Int64Histogram metric internally with custom buckets
- Bucket boundaries: [5000, 10000, 20000, 50000, 100000, 200000, 400000, 600000, 800000, 1000000]
- Metric name: feeds.text_extractor.body_char_count

Acceptance Criteria:
- Text extractor respects configurable limit
- Logs warning at truncation with appropriate details
- Emits histogram metric for all extractions

Context & Research:
- See internal/html/text_extractor.go for current implementation
- See internal/otel.go for how to get meter (otel.Meter(serviceName))
- Use metric.WithExplicitBucketBoundaries for custom buckets

Dependencies: Task 1 (config added)
