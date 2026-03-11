package textextractor

import (
	"context"
	"strings"
	"testing"
)

func TestExtractArticleText_SimpleHTML(t *testing.T) {
	html := `<html><body><h1>Title</h1><p>Paragraph one.</p><p>Paragraph two.</p></body></html>`
	extractor := ExtractArticleText(10000)
	text, ok := extractor(context.Background(), html)
	if !ok {
		t.Fatalf("expected success extracting text, got ok=false")
	}
	if len(text) == 0 {
		t.Fatalf("expected non-empty extracted text")
	}
	if !strings.Contains(text, "Paragraph one.") || !strings.Contains(text, "Paragraph two.") {
		t.Fatalf("extracted text does not contain expected content: %q", text)
	}
}

func TestExtractArticleText_Truncation(t *testing.T) {
	html := `<html><body><p>This is a long paragraph that should be truncated when the limit is low.</p></body></html>`
	extractor := ExtractArticleText(10)
	text, ok := extractor(context.Background(), html)
	if !ok {
		t.Fatalf("expected success extracting text even with truncation, got ok=false")
	}
	if len(text) != 10 {
		t.Fatalf("expected text length 10, got %d: %q", len(text), text)
	}
}

func TestStripHTMLToPlainText(t *testing.T) {
	html := `<html><body><script>bad</script><p>Hello <b>World</b></p></body></html>`
	plain := StripHTMLToPlainText(html)
	if len(plain) == 0 {
		t.Fatalf("expected non-empty plain text")
	}
	if !strings.Contains(plain, "Hello") || !strings.Contains(plain, "World") {
		t.Fatalf("expected stripped text to contain words, got: %q", plain)
	}
}
