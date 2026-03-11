package activity

import (
	"context"
	"encoding/json"
	"errors"
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

// ErrInvalidCategoryCount is returned when the number of categories is not between 1 and 5
var ErrInvalidCategoryCount = errors.New("categories must be between 1 and 5")

// CategorizeContent reads the fetched HTML content, sends it to the LLM for categorization,
// and saves the categories back to the MongoDB document.
//
// @param ctx - Context for the activity
// @param feedItemDoc - The FeedItemDocument containing the document ID, link, and title
// @return error - Returns an error if categorization fails or JSON parsing fails
// @author Thomas De Meyer
func CategorizeContent(c *mongo.Collection, client openai.Client, model, dataDir string) func(ctx context.Context, feedItemDoc internal.FeedItemDocument) error {
	return func(ctx context.Context, feedItemDoc internal.FeedItemDocument) error {
		logger := activity.GetLogger(ctx)

		filename := filepath.Join(dataDir, feedItemDoc.ID.Hex()+".html")
		htmlContent, err := os.ReadFile(filename)
		if err != nil {
			logger.Error("Failed to read HTML file for categorization", "err", err, "filename", filename)
			return err
		}

		logger.Info("Read HTML file for categorization", "id", feedItemDoc.ID.Hex(), "size", len(htmlContent))

		// Extract article text
		articleText, ok := textextractor.ExtractArticleText(string(htmlContent))
		if !ok || len(articleText) == 0 {
			articleText = textextractor.StripHTMLToPlainText(string(htmlContent))
		}

		// Build categorization prompt
		promptText := prompt.BuildCategorizationPrompt(feedItemDoc.Title, feedItemDoc.Link, articleText)

		// Call LLM for categorization
		resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
			Model: openai.ChatModel(model),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(promptText),
			},
		})
		if err != nil {
			logger.Error("Failed to generate categories with LLM", "err", err, "id", feedItemDoc.ID.Hex())
			return err
		}

		categoriesResponse := resp.Choices[0].Message.Content

		logger.Info("Received categories response", "id", feedItemDoc.ID.Hex(), "responseLength", len(categoriesResponse))

		// Parse JSON response to extract categories
		var categories []string
		if err := json.Unmarshal([]byte(categoriesResponse), &categories); err != nil {
			logger.Error("Failed to parse categories JSON", "err", err, "id", feedItemDoc.ID.Hex(), "response", categoriesResponse)
			return err
		}

		// Validate we have 1-5 categories
		if len(categories) == 0 || len(categories) > 5 {
			logger.Error("Invalid number of categories", "count", len(categories), "id", feedItemDoc.ID.Hex())
			return ErrInvalidCategoryCount
		}

		logger.Info("Parsed categories", "id", feedItemDoc.ID.Hex(), "categories", categories)

		// Update MongoDB document with categories
		filter := bson.M{"_id": feedItemDoc.ID}
		update := bson.M{"$set": bson.M{"categories": categories}}

		_, err = c.UpdateOne(ctx, filter, update)
		if err != nil {
			logger.Error("Failed to update document with categories", "err", err, "id", feedItemDoc.ID.Hex())
			return err
		}

		logger.Info("Successfully saved categories to document", "id", feedItemDoc.ID.Hex(), "link", feedItemDoc.Link, "categories", categories)
		return nil
	}
}
