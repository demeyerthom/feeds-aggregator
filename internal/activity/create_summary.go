package activity

import (
	"context"
	"os"
	"path/filepath"

	"github.com/demeyerthom/feeds-aggregator/internal"
	textextractor "github.com/demeyerthom/feeds-aggregator/internal/html"
	prompt "github.com/demeyerthom/feeds-aggregator/internal/prompt"
	"github.com/ollama/ollama/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.temporal.io/sdk/activity"
)

// CreateSummary reads the fetched HTML content, sends it to Ollama for summarization,
// and saves the summary back to the MongoDB document.
//
// @param ctx - Context for the activity
// @param feedItemDoc - The FeedItemDocument containing the document ID and link
// @return error - Returns an error if summary creation fails
// @author GitHub Copilot
func CreateSummary(c *mongo.Collection, ol *api.Client, model, dataDir string) func(ctx context.Context, feedItemDoc internal.FeedItemDocument) error {
	return func(ctx context.Context, feedItemDoc internal.FeedItemDocument) error {
		logger := activity.GetLogger(ctx)

		filename := filepath.Join(dataDir, feedItemDoc.ID.Hex()+".html")
		htmlContent, err := os.ReadFile(filename)
		if err != nil {
			logger.Error("Failed to read HTML file", "err", err, "filename", filename)
			return err
		}

		logger.Info("Read HTML file", "id", feedItemDoc.ID.Hex(), "size", len(htmlContent))

		// Extract article text and build a robust prompt
		articleText, ok := textextractor.ExtractArticleText(string(htmlContent))
		if !ok || len(articleText) == 0 {
			articleText = textextractor.StripHTMLToPlainText(string(htmlContent))
		}
		promptText := prompt.BuildSummaryPrompt("", feedItemDoc.Link, articleText)

		var summary string
		req := &api.GenerateRequest{
			Model:  model,
			Prompt: promptText,
		}

		err = ol.Generate(ctx, req, func(resp api.GenerateResponse) error {
			summary += resp.Response
			return nil
		})
		if err != nil {
			logger.Error("Failed to generate summary with Ollama", "err", err, "id", feedItemDoc.ID.Hex())
			return err
		}

		logger.Info("Generated summary", "id", feedItemDoc.ID.Hex(), "summaryLength", len(summary))

		filter := bson.M{"_id": feedItemDoc.ID}
		update := bson.M{"$set": bson.M{"summary": summary}}

		_, err = c.UpdateOne(ctx, filter, update)
		if err != nil {
			logger.Error("Failed to update document with summary", "err", err, "id", feedItemDoc.ID.Hex())
			return err
		}

		logger.Info("Successfully saved summary to document", "id", feedItemDoc.ID.Hex(), "link", feedItemDoc.Link)
		return nil
	}
}
