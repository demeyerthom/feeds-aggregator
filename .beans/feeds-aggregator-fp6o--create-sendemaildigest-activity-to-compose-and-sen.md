---
# feeds-aggregator-fp6o
title: Create SendEmailDigest activity to compose and send the daily email
status: todo
type: task
priority: normal
created_at: 2026-02-28T20:14:59Z
updated_at: 2026-02-28T20:41:29Z
parent: feeds-aggregator-hk58
blocked_by:
    - feeds-aggregator-hrta
---

## Context

Once the top N articles are selected, we need to compose a readable email and deliver it. No email capability exists in the project today. This is a brand-new activity.

## What to do

1. **Create `internal/activity/send_email_digest.go`** following the closure pattern. The outer function accepts SMTP configuration (or an `*smtp.Client` / mail-sending interface); the inner function is the Temporal activity:
   ```go
   func SendEmailDigest(smtpCfg SmtpConfig) func(ctx context.Context, input EmailDigestInput) error
   ```
2. **Define types** (in `internal/types.go`):
   - `SmtpConfig` struct — `Host`, `Port`, `Username`, `Password`, `From` (sender address)
   - `EmailDigestInput` struct — `To string` (recipient), `Articles []FeedItemDocument`, `Date time.Time`
3. **Compose the email body:**
   - **Subject:** `"Your Daily News Digest — {date}"`
   - **HTML body:** A clean, minimal HTML email listing each article as:
     - Title (as a clickable link to `article.Link`)
     - Summary (the 2–3 sentence Ollama summary)
     - Source / published date if available
   - **Plain-text fallback** for clients that don't render HTML
   - Use Go's `html/template` for the HTML body — define the template as an embedded string (`embed` or const).
4. **Send via SMTP** using Go's `net/smtp` package (stdlib — no extra dependency needed):
   - Connect to `SmtpConfig.Host:Port`
   - Authenticate with `smtp.PlainAuth`
   - Send the composed MIME multipart message (HTML + plain text)
5. **Error handling:**
   - SMTP connection/auth failures → return error (workflow will retry)
   - Empty article list → return early with a log message, no email sent, no error
6. **Security:**
   - Credentials come from env vars, never hardcoded
   - Validate `To` address format before attempting send

## Files to create / modify

- **New:** `internal/activity/send_email_digest.go`
- **New:** `internal/activity/send_email_digest_test.go`
- **Modify:** `internal/types.go` — add `SmtpConfig`, `EmailDigestInput`

## Acceptance criteria

- [ ] New `SendEmailDigest` activity follows the closure pattern
- [ ] `SmtpConfig` and `EmailDigestInput` types defined
- [ ] HTML email template renders article list with clickable titles, summaries
- [ ] Plain-text fallback included in MIME message
- [ ] Sends via `net/smtp` with `PlainAuth`
- [ ] Empty article list → no email sent, no error
- [ ] Invalid recipient address → returns descriptive error
- [ ] Unit tests cover: template rendering, empty list, SMTP error propagation (mock the SMTP send)
