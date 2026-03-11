---
# feeds-aggregator-1zrq
title: Add TEXT_LIMIT config to worker
status: scrapped
type: task
priority: normal
created_at: 2026-03-11T17:33:13Z
updated_at: 2026-03-11T17:35:00Z
parent: feeds-aggregator-vm05
---

Description: Add TEXT_LIMIT environment variable configuration to worker main.go with default 400000.\n\nOutput Requirements: Config struct with env:"TEXT_LIMIT,default=400000"\n\nAcceptance Criteria: Config loads correctly from environment with default value\n\nContext & Research:\n- See cmd/worker/main.go lines 40-65 for existing config patterns using go-env\n- The config should follow the same pattern as other config fields\n\nDependencies: None (this task can start immediately)

## Reasons for Scrapping\n\nDuplicate - using feeds-aggregator-e8t3 instead
