package prompt

import (
	"strings"
	"testing"
)

func TestBuildCategorizationPrompt_Basic(t *testing.T) {
	title := "Sample Title"
	url := "https://example.com/article"
	text := "This is a test content extracted from article."
	p := BuildCategorizationPrompt(title, url, text)

	if len(p) == 0 {
		t.Fatalf("prompt should not be empty")
	}
	if !strings.Contains(p, title) || !strings.Contains(p, url) || !strings.Contains(p, text) {
		t.Fatalf("prompt does not include input data: %q", p)
	}
}

func TestBuildCategorizationPrompt_ContainsExemplars(t *testing.T) {
	p := BuildCategorizationPrompt("Title", "https://example.com", "Content")

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

func TestBuildCategorizationPrompt_InstructionToCreateCategories(t *testing.T) {
	p := BuildCategorizationPrompt("Title", "https://example.com", "Content")

	if !strings.Contains(p, "encouraged to create your own categories") {
		t.Error("prompt should instruct LLM it can create its own categories")
	}
}

func TestBuildCategorizationPrompt_JSONArrayInstruction(t *testing.T) {
	p := BuildCategorizationPrompt("Title", "https://example.com", "Content")

	if !strings.Contains(p, "valid JSON array") {
		t.Error("prompt should instruct LLM to return valid JSON array")
	}
	if !strings.Contains(p, "1-5") {
		t.Error("prompt should specify 1-5 categories")
	}
	if !strings.Contains(p, "No preamble") {
		t.Error("prompt should instruct no preamble")
	}
}
