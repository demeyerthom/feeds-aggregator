// Package activity provides Temporal activity implementations for the feeds aggregator.
package activity

import (
	"context"
	"errors"
	"time"

	"github.com/demeyerthom/feeds-aggregator/internal"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.temporal.io/sdk/activity"
)

// AddNewFeedItem inserts a new feed item into the repository.
//
// @param ctx - Context for the activity
// @param feedItem - The feed item to insert
// @return FeedItemDocument - The inserted document with ID
// @return error - Returns an error if insertion fails
// @author GitHub Copilot
func AddNewFeedItem(c *mongo.Collection) func(ctx context.Context, feedItem internal.FeedItem) (internal.FeedItemDocument, error) {
	return func(ctx context.Context, feedItem internal.FeedItem) (internal.FeedItemDocument, error) {

		logger := activity.GetLogger(ctx)

		doc := internal.FeedItemDocument{
			Link:      feedItem.Link,
			Title:     feedItem.Title,
			CreatedAt: time.Now(),
		}

		result, err := c.InsertOne(ctx, doc)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				logger.Info("Feed item already exists", "link", feedItem.Link)
				return internal.FeedItemDocument{}, errors.New("feed item already exists")
			}
			logger.Error("Failed to insert feed item", "err", err)
			return internal.FeedItemDocument{}, err
		}

		doc.ID = result.InsertedID.(primitive.ObjectID)

		logger.Info("Successfully inserted feed item", "id", doc.ID.Hex(), "link", feedItem.Link, "title", feedItem.Title)
		return doc, nil
	}
}
