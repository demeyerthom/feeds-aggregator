package prompt

import "fmt"

// BuildSummaryPrompt constructs the prompt for the summarization model.
func BuildSummaryPrompt(title, url, text string) string {
	// A compact, instruction-focused prompt for news article summarization.
	// The instruction explicitly forbids preamble so the model outputs only the summary text.
	return fmt.Sprintf("You are a news summarizer. Output only a 2-3 sentence summary of the key facts — no preamble, no titles, no meta-commentary. Title: %s; URL: %s; Content: %s", title, url, text)
}
