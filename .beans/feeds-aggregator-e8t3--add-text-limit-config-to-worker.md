---
# feeds-aggregator-e8t3
title: Add TEXT_LIMIT config to worker
status: completed
type: task
priority: normal
created_at: 2026-03-11T17:33:17Z
updated_at: 2026-03-11T19:11:36Z
parent: feeds-aggregator-vm05
---

Description: Add TEXT_LIMIT environment variable configuration to worker main.go with default 400000.

Output Requirements: Config struct with TextLimit field using go-env env tag

Acceptance Criteria: Config loads correctly from environment with default value

Context & Research:
- See cmd/worker/main.go lines 40-65 for existing config patterns using go-env
- The config should follow the same pattern as other config fields

Dependencies: None

## Coder Notes\n- Added TextExtractor struct with Limit field to Configuration\n- Environment variable: TEXT_LIMIT\n- Default value: 400000
