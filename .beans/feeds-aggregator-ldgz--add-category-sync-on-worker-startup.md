---
# feeds-aggregator-ldgz
title: Add category sync on worker startup
status: todo
type: task
created_at: 2026-03-11T20:37:44Z
updated_at: 2026-03-11T20:37:44Z
parent: feeds-aggregator-4kxe
---

Description: On worker startup, sync existing categories from MongoDB to Redis. Use smart sync pattern: only sync if Redis set is empty (SCARD returns 0). This handles multiple workers gracefully - first worker to start does the sync, others skip.

Output Requirements:
- Add sync logic to cmd/worker/main.go after Redis and MongoDB initialization
- Check Redis SCARD feeds:categories
- If count == 0: query MongoDB distinct("categories") and SADD to feeds:categories
- Log sync status (synced count or skipped)

Acceptance Criteria:
- go build ./cmd/worker succeeds
- Worker syncs categories on startup if Redis set is empty
- Worker skips sync if Redis set already has categories
- Categories are stored in Redis Set feeds:categories

Context & Research:
- Redis Set operations: SCARD (count), SADD (add members)
- MongoDB distinct query: collection.Distinct(ctx, "categories", bson.M{})
- Add after MongoDB connection is established (around line 168)
- Use rdb.SCard and rdb.SAdd from go-redis

Open Questions: None

Dependencies: None
