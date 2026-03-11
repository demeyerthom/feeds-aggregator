---
# feeds-aggregator-4kxe
title: Use Redis-cached known categories for smarter categorization
status: todo
type: feature
created_at: 2026-03-11T20:37:39Z
updated_at: 2026-03-11T20:37:39Z
parent: feeds-aggregator-vxnm
---

Description: Cache known categories in Redis and include them in the LLM prompt to encourage reuse of existing categories over creating new ones.

General Requirements:
- Redis Set `feeds:categories` stores all known categories
- On worker startup: sync categories from MongoDB to Redis (if empty)
- ProcessContent activity fetches known categories from Redis
- Prompt includes both taxonomy exemplars AND known categories
- New categories are added to Redis Set immediately
- Handles multiple workers scaling gracefully

Design Choices:
- Redis key: `feeds:categories` (namespaced)
- Smart sync: Only sync from MongoDB if Redis set is empty
- Prompt structure: Taxonomy exemplars + Known categories (preferred)
- Idempotent operations for scaling safety
