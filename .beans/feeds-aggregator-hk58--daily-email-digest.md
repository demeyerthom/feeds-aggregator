---
# feeds-aggregator-hk58
title: Daily email digest
status: todo
type: feature
priority: normal
created_at: 2026-02-28T20:41:23Z
updated_at: 2026-02-28T20:59:03Z
parent: feeds-aggregator-65wn
blocked_by:
    - feeds-aggregator-9wbo
---

Build and wire everything needed to select the top articles from MongoDB and send a daily digest email: SelectTopArticles activity, SendEmailDigest activity, DailyDigest workflow, worker registration, and cron scheduling.

## Summary of Changes\n- Closed the Daily Digest feature. Implemented tasks: SelectTopArticles, SendEmailDigest, DailyDigest workflow, worker registration, and cron scheduling. All related changes are now completed and ready for review.
