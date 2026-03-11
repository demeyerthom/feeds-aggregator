package prompt

import "fmt"

// BuildCategorizationPrompt constructs the prompt for categorizing content.
// It includes the taxonomy as exemplars but instructs the LLM it can create its own categories.
// Returns a JSON array of 1-5 categories.
func BuildCategorizationPrompt(title, url, text string) string {
	// EXEMPLARS show the taxonomy structure but the LLM is instructed to create its own categories.
	return fmt.Sprintf(`You are a content categorizer. Analyze the following article and assign 1-5 categories.

EXEMPLARS (you may use these or create your own categories):
1. Programming Languages - Language releases, comparisons, best practices, performance tips, ecosystem libraries. Examples: JavaScript, Rust, Go, Python
2. Frameworks & Libraries - Framework introductions, major updates, ecosystem tools, architecture patterns
3. Developer Tools - IDEs, CLI tools, build tools, debugging, testing, productivity
4. AI & Machine Learning - AI tools, ML concepts, LLMs, prompt engineering, AI frameworks
5. Software Engineering Practices - Architecture, code quality, testing, design patterns, refactoring
6. DevOps & Infrastructure - CI/CD, containers, cloud, observability, IaC, scaling
7. Security - Secure coding, vulnerabilities, authentication, cryptography, security tools
8. Industry & Trends - Tech trends, ecosystem shifts, market changes, startup stacks
9. Opinion & Thought Pieces - Opinions, predictions, lessons learned, DX, career
10. Tutorials & Guides - Beginner tutorials, step-by-step guides, walkthroughs, project builds

IMPORTANT: You are encouraged to create your own categories if the content doesn't fit the exemplars above. Think about what categories would best describe this content.

Return a valid JSON array of 1-5 category strings. No preamble, just the JSON array.

Title: %s
URL: %s
Content: %s`, title, url, text)
}
