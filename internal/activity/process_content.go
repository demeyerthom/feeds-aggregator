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

// processContentResponse represents the JSON response from the LLM
type processContentResponse struct {
	Summary    string   `json:"summary"`
	Categories []string `json:"categories"`
}

// ProcessContent reads the fetched HTML content, sends it to the LLM for combined
// summarization and categorization, and saves both to the MongoDB document in a single operation.
//
// @param c - MongoDB collection for updating feed item documents
// @param client - OpenAI client for LLM calls
// @param model - The model to use for LLM calls
// @param dataDir - Directory where HTML files are stored
// @param textLimit - Maximum characters to extract from HTML content
// @return A function that processes a feed item document
// @author Thomas De Meyer
func ProcessContent(c *mongo.Collection, client openai.Client, model, dataDir string, textLimit int) func(ctx context.Context, feedItemDoc internal.FeedItemDocument) error {
	return func(ctx context.Context, feedItemDoc internal.FeedItemDocument) error {
		logger := activity.GetLogger(ctx)

		filename := filepath.Join(dataDir, feedItemDoc.ID.Hex()+".html")
		htmlContent, err := os.ReadFile(filename)
		if err != nil {
			logger.Error("Failed to read HTML file for content processing", "err", err, "filename", filename)
			return err
		}

		logger.Info("Read HTML file for content processing", "id", feedItemDoc.ID.Hex(), "size", len(htmlContent))

		// Extract article text
		extractText := textextractor.ExtractArticleText(textLimit)
		articleText, ok := extractText(ctx, string(htmlContent))
		if !ok || len(articleText) == 0 {
			articleText = textextractor.StripHTMLToPlainText(string(htmlContent))
		}

		// Build combined prompt for summarization and categorization
		promptText := prompt.BuildProcessContentPrompt(feedItemDoc.Title, feedItemDoc.Link, articleText)

		// Call LLM for combined processing
		resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
			Model: openai.ChatModel(model),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(promptText),
			},
		})
		if err != nil {
			logger.Error("Failed to process content with LLM", "err", err, "id", feedItemDoc.ID.Hex())
			return err
		}

		llmResponse := resp.Choices[0].Message.Content

		logger.Info("Received LLM response", "id", feedItemDoc.ID.Hex(), "responseLength", len(llmResponse))

		// Parse JSON response to extract summary and categories
		var result processContentResponse
		if err := json.Unmarshal([]byte(llmResponse), &result); err != nil {
			logger.Error("Failed to parse LLM JSON response", "err", err, "id", feedItemDoc.ID.Hex(), "response", llmResponse)
			return err
		}

		// Validate we have 1-5 categories
		if len(result.Categories) == 0 || len(result.Categories) > 5 {
			logger.Error("Invalid number of categories", "count", len(result.Categories), "id", feedItemDoc.ID.Hex())
			return ErrInvalidCategoryCount
		}

		logger.Info("Parsed content processing result", "id", feedItemDoc.ID.Hex(), "summaryLength", len(result.Summary), "categories", result.Categories)

		// Update MongoDB document with both summary and categories in a single operation
		filter := bson.M{"_id": feedItemDoc.ID}
		update := bson.M{"$set": bson.M{"summary": result.Summary, "categories": result.Categories}}

		_, err = c.UpdateOne(ctx, filter, update)
		if err != nil {
			logger.Error("Failed to update document with summary and categories", "err", err, "id", feedItemDoc.ID.Hex())
			return err
		}

		logger.Info("Successfully saved summary and categories to document", "id", feedItemDoc.ID.Hex(), "link", feedItemDoc.Link, "categories", result.Categories)
		return nil
	}
}
