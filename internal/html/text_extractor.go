package textextractor

import (
	"bytes"
	"context"
	"log/slog"
	"regexp"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"golang.org/x/net/html"
)

var textLengthHistogram metric.Int64Histogram

func init() {
	meter := otel.Meter("feeds-worker")

	buckets := []float64{5000, 10000, 20000, 50000, 100000, 200000, 400000, 600000, 800000, 1000000}

	textLengthHistogram, _ = meter.Int64Histogram(
		"feeds.text_extractor.body_char_count",
		metric.WithDescription("Length of extracted article text"),
		metric.WithExplicitBucketBoundaries(buckets...),
	)
}

// ExtractArticleText creates an extractor function that parses HTML and extracts visible article text,
// stripping boilerplate like scripts, styles, navigation, headers, and footers.
// The limit parameter controls the maximum number of characters to extract.
// It returns a function that accepts a context and HTML string, and returns the extracted text and a success boolean.
func ExtractArticleText(limit int) func(ctx context.Context, htmlStr string) (string, bool) {
	return func(ctx context.Context, htmlStr string) (string, bool) {
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

		// Record metric for all extractions
		textLengthHistogram.Record(ctx, int64(len(text)))

		if text == "" {
			return "", false
		}
		if len(text) > limit {
			slog.Warn("Text truncated", "originalLength", len(text), "limit", limit)
			text = text[:limit]
		}
		return text, true
	}
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
