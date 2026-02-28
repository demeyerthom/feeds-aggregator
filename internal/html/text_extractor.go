package textextractor

import (
	"bytes"
	"golang.org/x/net/html"
	"regexp"
	"strings"
)

// ExtractArticleText parses HTML and extracts the visible article text,
// stripping boilerplate like scripts, styles, navigation, headers, and footers.
// It returns the extracted text (up to 4000 chars) and a boolean indicating success.
func ExtractArticleText(htmlStr string) (string, bool) {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return "", false
	}
	var buf bytes.Buffer
	var dfs func(n *html.Node, skip bool)
	isSkip := func(tag string) bool {
		switch tag {
		case "script", "style", "nav", "header", "footer", "aside":
			return true
		default:
			return false
		}
	}
	dfs = func(n *html.Node, skip bool) {
		if n.Type == html.ElementNode {
			if isSkip(n.Data) {
				skip = true
			}
		}
		if n.Type == html.TextNode && !skip {
			t := strings.TrimSpace(n.Data)
			if t != "" {
				buf.WriteString(t)
				buf.WriteString(" ")
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			dfs(c, skip)
		}
	}
	dfs(doc, false)
	text := strings.TrimSpace(buf.String())
	if text == "" {
		return "", false
	}
	if len(text) > 4000 {
		text = text[:4000]
	}
	return text, true
}

// StripHTMLToPlainText is a resilient fallback that removes HTML tags to plain text.
func StripHTMLToPlainText(htmlStr string) string {
	// Simple tag stripper as a fallback (non-boilerplate aware)
	re := regexp.MustCompile("<[^>]+>")
	text := re.ReplaceAllString(htmlStr, " ")
	text = strings.TrimSpace(text)
	// Normalize whitespace
	fields := strings.Fields(text)
	return strings.Join(fields, " ")
}
