package activity

import (
	"context"
	"os"
	"path/filepath"

	"github.com/demeyerthom/feeds-aggregator/internal"
	textextractor "github.com/demeyerthom/feeds-aggregator/internal/html"
	prompt "github.com/demeyerthom/feeds-aggregator/internal/prompt"
	"github.com/openai/openai-go/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.temporal.io/sdk/activity"
)

// CreateSummary reads the fetched HTML content, sends it to Zen for summarization,
// and saves the summary back to the MongoDB document.
//
// @param ctx - Context for the activity
// @param feedItemDoc - The FeedItemDocument containing the document ID and link
// @return error - Returns an error if summary creation fails
// @author GitHub Copilot
func CreateSummary(c *mongo.Collection, client openai.Client, model, dataDir string, textLimit int) func(ctx context.Context, feedItemDoc internal.FeedItemDocument) error {
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
		extractText := textextractor.ExtractArticleText(textLimit)
		articleText, ok := extractText(string(htmlContent))
		if !ok || len(articleText) == 0 {
			articleText = textextractor.StripHTMLToPlainText(string(htmlContent))
		}
		promptText := prompt.BuildSummaryPrompt("", feedItemDoc.Link, articleText)

		resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
			Model: openai.ChatModel(model),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(promptText),
			},
		})
		if err != nil {
			logger.Error("Failed to generate summary with Zen", "err", err, "id", feedItemDoc.ID.Hex())
			return err
		}

		summary := resp.Choices[0].Message.Content

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
