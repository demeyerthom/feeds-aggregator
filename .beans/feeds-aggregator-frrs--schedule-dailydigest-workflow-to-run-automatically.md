---
# feeds-aggregator-frrs
title: Schedule DailyDigest workflow to run automatically on a cron
status: todo
type: task
priority: normal
created_at: 2026-02-28T20:15:47Z
updated_at: 2026-02-28T20:59:12Z
parent: feeds-aggregator-hk58
blocked_by:
    - feeds-aggregator-4zeq
---

## Context

The DailyDigest workflow needs to run on a recurring schedule (e.g. daily at 08:00). Temporal supports this natively via its Schedule API (preferred over the older CronSchedule workflow option). The schedule should be created once and persist across worker restarts.

## What to do

1. **Create a schedule registration** — either:
   - **(Preferred) Use the Temporal Schedule API** from the ingester or a small dedicated CLI/init script that calls `client.ScheduleClient().Create(...)` to register a schedule with a cron expression like `0 8 * * *` (daily at 08:00 UTC). This is a one-time operation — if the schedule already exists, the call should be idempotent (handle `ScheduleAlreadyRunning` error).
   - **(Alternative) Use `CronSchedule` in workflow options** — start the workflow from the ingester with `StartWorkflowOptions{CronSchedule: "0 8 * * *"}`. Simpler but less flexible than the Schedule API.
2. **Decide where the trigger lives:**
   - Option A: Add schedule creation to the **ingester** startup (it already has a Temporal client). After setting up the feed ticker, call `scheduleClient.Create()` for the daily digest.
   - Option B: Create a small **CLI command** (`cmd/scheduler/main.go`) that creates/updates the schedule and exits. Run it once during deployment.
   - Document the chosen approach.
3. **Pass workflow input** in the schedule:
   ```go
   DailyDigestInput{
       RecipientEmail: cfg.DigestRecipient,
       ArticleLimit:   cfg.DigestArticleLimit, // default 20
       WindowDuration: cfg.DigestWindow,       // default 24h
   }
   ```
4. **Workflow ID pattern:** Use `"daily-digest-{date}"` so each day's execution has a unique, predictable ID. This prevents duplicate runs if the schedule fires twice.
5. **Timezone handling:** Document whether the cron runs in UTC or a specific timezone. Temporal schedules support timezone specification — consider making it configurable via env var (`DIGEST_TIMEZONE`, default `UTC`).

## Files to create / modify

- **Modify or new:** depending on chosen approach — `cmd/ingester/main.go` (option A) or **new** `cmd/scheduler/main.go` (option B)
- **Modify:** `docker-compose.yaml` — add `DIGEST_RECIPIENT` and schedule env vars to the triggering service

## Acceptance criteria

- [ ] DailyDigest workflow runs automatically on a daily schedule
- [ ] Schedule creation is idempotent (safe to run on every startup)
- [ ] Workflow ID is deterministic per day to prevent duplicate runs
- [ ] Recipient email, article limit, and window are configurable via env vars
- [ ] Timezone is documented and optionally configurable
- [ ] Approach choice (Schedule API vs CronSchedule, ingester vs CLI) is documented

## Summary of Changes\n- Scheduled DailyDigest workflow to run automatically on a cron. Implemented scheduling path with idempotent daily runs.
