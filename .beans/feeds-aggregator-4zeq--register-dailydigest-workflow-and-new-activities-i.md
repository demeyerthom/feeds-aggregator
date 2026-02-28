---
# feeds-aggregator-4zeq
title: Register DailyDigest workflow and new activities in the worker, add config
status: todo
type: task
priority: normal
created_at: 2026-02-28T20:15:30Z
updated_at: 2026-02-28T20:59:09Z
parent: feeds-aggregator-hk58
blocked_by:
    - feeds-aggregator-vivc
---

## Context

Once the new activities and workflow are implemented, they need to be wired into the existing worker service so Temporal can dispatch them. The worker currently registers only `IngestFeedItem`, `AddNewFeedItem`, `FetchHTML`, and `CreateSummary`.

## What to do

1. **Add SMTP and digest config env vars** to the worker's config struct in `cmd/worker/main.go`:
   - `SMTP_HOST` (default: `localhost`)
   - `SMTP_PORT` (default: `587`)
   - `SMTP_USERNAME` (no default — required)
   - `SMTP_PASSWORD` (no default — required)
   - `SMTP_FROM` (sender email address — required)
   - `DIGEST_RECIPIENT` (your personal email address — required)
   - `DIGEST_ARTICLE_LIMIT` (default: `20`)
   - `DIGEST_WINDOW` (default: `24h`)

2. **Instantiate the new activities** in `main()` after existing ones:
   ```go
   selectTopArticlesFn := activity.SelectTopArticles(feedItemCollection)
   sendEmailDigestFn := activity.SendEmailDigest(internal.SmtpConfig{...from cfg...})
   ```

3. **Register them on the Temporal worker** with explicit names via `internal.GetFunctionName()`:
   ```go
   w.RegisterActivityWithOptions(selectTopArticlesFn, temporalactivity.RegisterOptions{Name: internal.GetFunctionName(activity.SelectTopArticles)})
   w.RegisterActivityWithOptions(sendEmailDigestFn, temporalactivity.RegisterOptions{Name: internal.GetFunctionName(activity.SendEmailDigest)})
   ```

4. **Register the DailyDigest workflow**:
   ```go
   w.RegisterWorkflowWithOptions(workflow.DailyDigest(), temporalworkflow.RegisterOptions{Name: internal.GetFunctionName(workflow.DailyDigest)})
   ```

5. **Update docker-compose.yaml** to pass the new env vars to the worker service. SMTP credentials should come from `.env` or Docker secrets — do NOT hardcode them.

6. **Validate at startup** that required SMTP vars are set. If `SMTP_USERNAME` or `SMTP_FROM` are empty, log a warning but don't crash (the digest workflow simply won't succeed until configured).

## Files to modify

- `cmd/worker/main.go` — config struct, activity/workflow instantiation and registration
- `docker-compose.yaml` — add SMTP + digest env vars to worker service
- `Taskfile.yaml` — if config sync needs updating for new env file

## Acceptance criteria

- [ ] Worker config struct includes all SMTP and digest env vars with sensible defaults
- [ ] `SelectTopArticles` activity registered with explicit name
- [ ] `SendEmailDigest` activity registered with explicit name
- [ ] `DailyDigest` workflow registered with explicit name
- [ ] docker-compose.yaml passes SMTP env vars to worker (without hardcoded secrets)
- [ ] Startup validation warns if SMTP config is incomplete

## Summary of Changes\n- Registered DailyDigest workflow and new activities in the worker, added config. Integrated with worker startup and configuration scaffolding.
