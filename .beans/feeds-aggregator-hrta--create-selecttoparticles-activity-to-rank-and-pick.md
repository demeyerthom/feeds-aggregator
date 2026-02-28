---
# feeds-aggregator-hrta
title: Create SelectTopArticles activity to rank and pick top N articles
status: todo
type: task
priority: normal
created_at: 2026-02-28T20:14:41Z
updated_at: 2026-02-28T20:41:29Z
parent: feeds-aggregator-hk58
blocked_by:
    - feeds-aggregator-n6q1
---

## Context

After articles are ingested, stored in MongoDB, and summarised, we need a way to select the best N articles for a digest email. This activity does not exist yet — it is a brand-new addition.

## What to do

1. **Create `internal/activity/select_top_articles.go`** following the project's closure-based DI pattern:
   ```go
   func SelectTopArticles(c *mongo.Collection) func(ctx context.Context, input SelectTopArticlesInput) ([]internal.FeedItemDocument, error)
   ```
2. **Define an input type** (in `internal/types.go` or inline):
   - `WindowDuration time.Duration` — how far back to look (e.g. 24h)
   - `Limit int` — how many articles to return (default 20)
3. **Query MongoDB** for `FeedItemDocument` records where:
   - `created_at` is within the window
   - `summary` is non-empty (article has been successfully summarised)
   - Sort by `created_at` descending (most recent first)
   - Limit to `Limit` results
4. **Keep the scoring simple for now** — pure recency. Add a `// TODO: add source-quality weighting` comment for future enhancement. The important thing is that the query is efficient and the activity is well-tested.
5. **Return** the slice of `FeedItemDocument` (including `_id`, `title`, `link`, `summary`, `created_at`).
6. **Handle edge cases:**
   - Fewer than `Limit` articles available → return whatever exists (even 0)
   - MongoDB query error → return the error so the workflow can retry

## Files to create / modify

- **New:** `internal/activity/select_top_articles.go`
- **Modify:** `internal/types.go` — add `SelectTopArticlesInput` struct
- **New:** `internal/activity/select_top_articles_test.go`

## Acceptance criteria

- [ ] New `SelectTopArticles` activity exists following the closure pattern
- [ ] `SelectTopArticlesInput` type defined with `WindowDuration` and `Limit` fields
- [ ] Queries MongoDB with time window filter, non-empty summary filter, sorted by recency
- [ ] Returns up to `Limit` fully populated `FeedItemDocument` records
- [ ] Gracefully handles < N results (returns partial slice, no error)
- [ ] Unit tests cover: normal case, empty results, MongoDB error propagation
