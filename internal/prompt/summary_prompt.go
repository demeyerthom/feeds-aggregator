package prompt

import "fmt"

// BuildSummaryPrompt constructs the prompt for the summarization model.
func BuildSummaryPrompt(title, url, text string) string {
	// A compact, instruction-focused prompt for news article summarization
	return fmt.Sprintf("Summarize the following news article in 2-3 sentences, focusing on the key facts. Title: %s; URL: %s; Content: %s", title, url, text)
}
