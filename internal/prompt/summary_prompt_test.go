package prompt

import "strings"
import "testing"

func TestBuildSummaryPrompt_Basic(t *testing.T) {
	title := "Sample Title"
	url := "https://example.com/article"
	text := "This is a test content extracted from article."
	p := BuildSummaryPrompt(title, url, text)
	if len(p) == 0 {
		t.Fatalf("prompt should not be empty")
	}
	if !strings.Contains(p, title) || !strings.Contains(p, url) || !strings.Contains(p, text) {
		t.Fatalf("prompt does not include input data: %q", p)
	}
}
