package prompt

import (
	"strings"
	"testing"
)

func TestBuildProcessContentPrompt_Basic(t *testing.T) {
	title := "Sample Title"
	url := "https://example.com/article"
	text := "This is a test content extracted from article."
	p := BuildProcessContentPrompt(title, url, text)

	if len(p) == 0 {
		t.Fatalf("prompt should not be empty")
	}
	if !strings.Contains(p, title) || !strings.Contains(p, url) || !strings.Contains(p, text) {
		t.Fatalf("prompt does not include input data: %q", p)
	}
}

func TestBuildProcessContentPrompt_ContainsExemplars(t *testing.T) {
	p := BuildProcessContentPrompt("Title", "https://example.com", "Content")

	exemplars := []string{
		"Programming Languages",
		"Frameworks & Libraries",
		"Developer Tools",
		"AI & Machine Learning",
		"Software Engineering Practices",
		"DevOps & Infrastructure",
		"Security",
		"Industry & Trends",
		"Opinion & Thought Pieces",
		"Tutorials & Guides",
	}

	for _, ex := range exemplars {
		if !strings.Contains(p, ex) {
			t.Errorf("prompt should contain exemplar %q", ex)
		}
	}
}

func TestBuildProcessContentPrompt_ContainsSummaryInstructions(t *testing.T) {
	p := BuildProcessContentPrompt("Title", "https://example.com", "Content")

	if !strings.Contains(p, "2-3 sentence summary") {
		t.Error("prompt should contain summary instructions for 2-3 sentences")
	}
	if !strings.Contains(p, "key facts") {
		t.Error("prompt should instruct to focus on key facts")
	}
}

func TestBuildProcessContentPrompt_ContainsCategorizationInstructions(t *testing.T) {
	p := BuildProcessContentPrompt("Title", "https://example.com", "Content")

	if !strings.Contains(p, "1-5 categories") {
		t.Error("prompt should specify 1-5 categories")
	}
	if !strings.Contains(p, "encouraged to create your own categories") {
		t.Error("prompt should instruct LLM it can create its own categories")
	}
}

func TestBuildProcessContentPrompt_JSONOutputInstruction(t *testing.T) {
	p := BuildProcessContentPrompt("Title", "https://example.com", "Content")

	if !strings.Contains(p, `"summary"`) {
		t.Error("prompt should instruct JSON output with summary field")
	}
	if !strings.Contains(p, `"categories"`) {
		t.Error("prompt should instruct JSON output with categories field")
	}
	if !strings.Contains(p, "No preamble") {
		t.Error("prompt should instruct no preamble")
	}
	if !strings.Contains(p, "valid JSON object") {
		t.Error("prompt should instruct valid JSON object output")
	}
}
