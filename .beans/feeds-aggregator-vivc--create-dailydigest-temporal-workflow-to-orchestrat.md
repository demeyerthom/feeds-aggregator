---
# feeds-aggregator-vivc
title: Create DailyDigest Temporal workflow to orchestrate select → email
status: todo
type: task
priority: normal
created_at: 2026-02-28T20:15:13Z
updated_at: 2026-02-28T20:59:06Z
parent: feeds-aggregator-hk58
blocked_by:
    - feeds-aggregator-fp6o
---

## Context

The existing `IngestFeedItem` workflow handles the per-article pipeline (store → fetch HTML → summarise). We need a **separate** workflow that runs on a schedule and orchestrates the digest: select the top articles, then send an email. This is a brand-new workflow.

## What to do

1. **Create `internal/workflow/daily_digest.go`** following the project's closure pattern:
   ```go
   func DailyDigest() func(ctx workflow.Context, input DailyDigestInput) error
   ```
2. **Define `DailyDigestInput`** (in `internal/types.go`):
   - `RecipientEmail string`
   - `ArticleLimit int` (default 20)
   - `WindowDuration time.Duration` (default 24h)
3. **Workflow steps** (sequential):
   1. Call `SelectTopArticles` activity with `{WindowDuration, ArticleLimit}` → receive `[]FeedItemDocument`
   2. If the slice is empty, log a warning and return early (no email to send) — this is **not** an error
   3. Call `SendEmailDigest` activity with `{To: RecipientEmail, Articles: results, Date: workflow.Now(ctx)}`
4. **Activity options:**
   - `SelectTopArticles`: `StartToCloseTimeout: 30s`, `MaximumAttempts: 3` (MongoDB query — should be fast, worth retrying on transient errors)
   - `SendEmailDigest`: `StartToCloseTimeout: 60s`, `MaximumAttempts: 3` (SMTP send — can retry on transient failures)
5. **Idempotency:** The workflow should be safe to re-run. Selecting and emailing the same articles twice is acceptable (it's a read-only query + email send), but use a deterministic workflow ID like `"daily-digest-{YYYY-MM-DD}"` so Temporal deduplicates same-day runs.

## Files to create / modify

- **New:** `internal/workflow/daily_digest.go`
- **Modify:** `internal/types.go` — add `DailyDigestInput`
- **New:** `internal/workflow/daily_digest_test.go`

## Acceptance criteria

- [ ] New `DailyDigest` workflow follows the closure pattern
- [ ] `DailyDigestInput` type defined with `RecipientEmail`, `ArticleLimit`, `WindowDuration`
- [ ] Workflow calls `SelectTopArticles` then `SendEmailDigest` sequentially
- [ ] Empty article list → logs warning, returns nil (no error, no email)
- [ ] Each activity has appropriate timeout and retry policy
- [ ] Unit tests cover: happy path (mock both activities), empty-result path, activity failure propagation

## Summary of Changes\n- Completed DailyDigest workflow implementation to orchestrate article selection and email delivery.
