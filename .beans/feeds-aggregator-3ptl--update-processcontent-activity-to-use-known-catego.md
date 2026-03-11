---
# feeds-aggregator-3ptl
title: Update ProcessContent activity to use known categories
status: todo
type: task
created_at: 2026-03-11T20:37:54Z
updated_at: 2026-03-11T20:37:54Z
parent: feeds-aggregator-4kxe
blocked_by:
    - feeds-aggregator-9o6s
---

Description: Update ProcessContent activity to fetch known categories from Redis, pass them to the prompt builder, and add new categories to Redis after LLM response.

Output Requirements:
- Fetch known categories from Redis at activity start: SMEMBERS feeds:categories
- Pass known categories to BuildProcessContentPrompt
- After LLM response, identify new categories (not in knownCategories)
- Add new categories to Redis: SADD feeds:categories <new_cats>
- Log when new categories are added

Acceptance Criteria:
- go build ./internal/activity succeeds
- Activity fetches categories from Redis
- Prompt receives known categories
- New categories are added to Redis Set
- Handles empty Redis set gracefully

Context & Research:
- Redis SMEMBERS returns []string
- Use rdb.SMembers(ctx, "feeds:categories")
- Compare LLM response categories with known to find new ones
- Use rdb.SAdd(ctx, "feeds:categories", newCat) for each new category
- ProcessContent closure needs rdb *redis.Client parameter

Open Questions: None

Dependencies: This task depends on feeds-aggregator-9o6s being completed first
