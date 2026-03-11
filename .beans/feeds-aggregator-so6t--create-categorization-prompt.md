---
# feeds-aggregator-so6t
title: Create categorization prompt
status: todo
type: task
priority: normal
created_at: 2026-03-11T16:26:24Z
updated_at: 2026-03-11T16:31:00Z
parent: feeds-aggregator-nfb7
---

Description: Create a new prompt file in internal/prompt/categorization_prompt.go that builds a prompt for the LLM to categorize content. The prompt should include the taxonomy (10 categories) as EXEMPLARS/EXAMPLES but instruct the LLM it is allowed to come up with its own categorizations. The LLM should return 1-5 categories as a JSON array. Create a corresponding test file.

## Category Taxonomy (use as examples, LLM can create its own)

1. **Programming Languages** - Language releases, comparisons, best practices, performance tips, ecosystem libraries. Examples: JavaScript, Rust, Go, Python
2. **Frameworks & Libraries** - Framework introductions, major updates, ecosystem tools, architecture patterns
3. **Developer Tools** - IDEs, CLI tools, build tools, debugging, testing, productivity
4. **AI & Machine Learning** - AI tools, ML concepts, LLMs, prompt engineering, AI frameworks
5. **Software Engineering Practices** - Architecture, code quality, testing, design patterns, refactoring
6. **DevOps & Infrastructure** - CI/CD, containers, cloud, observability, IaC, scaling
7. **Security** - Secure coding, vulnerabilities, authentication, cryptography, security tools
8. **Industry & Trends** - Tech trends, ecosystem shifts, market changes, startup stacks
9. **Opinion & Thought Pieces** - Opinions, predictions, lessons learned, DX, career
10. **Tutorials & Guides** - Beginner tutorials, step-by-step guides, walkthroughs, project builds

## Alternative (also valid)
- Topics: AI, Programming Languages, DevOps, Security, Tools
- Post Types: Tutorial, Opinion, Release Update, Deep Dive, Comparison

Output Requirements:
- New file internal/prompt/categorization_prompt.go with BuildCategorizationPrompt function
- New file internal/prompt/categorization_prompt_test.go with tests
- Prompt includes the category taxonomy as EXEMPLARS
- Prompt explicitly states LLM can create its own categories

Acceptance Criteria:
- Prompt function exists and returns properly formatted prompt string
- Tests pass
- Prompt instructs LLM to return valid JSON array
