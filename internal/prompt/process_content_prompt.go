package prompt

import "fmt"

// BuildProcessContentPrompt constructs a combined prompt for summarization and categorization.
// The prompt instructs the LLM to return a JSON object with both summary and categories fields.
// @param title - The title of the content to process
// @param url - The URL of the content
// @param text - The content text to summarize and categorize
// @return A prompt string that requests JSON output with summary and categories
// @author feeds-aggregator
func BuildProcessContentPrompt(title, url, text string) string {
	return fmt.Sprintf(`You are a content processor. Analyze the following article and provide both a summary and categories.

SUMMARY INSTRUCTIONS:
- Output a 2-3 sentence summary of the key facts
- No preamble, no titles, no meta-commentary
- Focus on the essential information

CATEGORIZATION INSTRUCTIONS:
- Assign 1-5 categories that best describe the content

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

OUTPUT FORMAT:
Return ONLY a valid JSON object with no preamble. The JSON must have this exact structure:
{"summary": "your 2-3 sentence summary here", "categories": ["category1", "category2"]}

Title: %s
URL: %s
Content: %s`, title, url, text)
}
