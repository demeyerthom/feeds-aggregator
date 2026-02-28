// Package workflow provides Temporal workflow implementations for the feeds aggregator.
package workflow

import (
	"time"

	"github.com/demeyerthom/feeds-aggregator/internal"
	"github.com/demeyerthom/feeds-aggregator/internal/activity"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// IngestFeedItem is the workflow function that orchestrates feed item ingestion.
// It executes three activities in sequence: add feed item, fetch HTML, and create summary.
//
// @param ctx - Workflow context
// @param feedItem - The feed item to ingest
// @return error - Returns an error if any activity fails
// @author GitHub Copilot
func IngestFeedItem() func(ctx workflow.Context, feedItem internal.FeedItem) error {
	return func(ctx workflow.Context, feedItem internal.FeedItem) error {
		workflow.GetLogger(ctx).Info("Ingest feed item workflow started.", "link", feedItem.Link, "title", feedItem.Title)

		ao := workflow.ActivityOptions{
			StartToCloseTimeout: 5 * time.Minute,
			RetryPolicy: &temporal.RetryPolicy{
				MaximumAttempts: 3,
			},
		}
		ctx = workflow.WithActivityOptions(ctx, ao)

		// First activity: add feed item to MongoDB and get the document with ID
		var feedItemDoc internal.FeedItemDocument
		err := workflow.ExecuteActivity(ctx, internal.GetFunctionName(activity.AddNewFeedItem), feedItem).Get(ctx, &feedItemDoc)
		if err != nil {
			workflow.GetLogger(ctx).Error("addNewFeedItemActivity activity failed.", "Error", err)
			return err
		}

		// Second activity: fetch HTML page and store to disk
		err = workflow.ExecuteActivity(ctx, internal.GetFunctionName(activity.FetchHTML), feedItemDoc).Get(ctx, nil)
		if err != nil {
			workflow.GetLogger(ctx).Error("fetchHTMLActivity activity failed.", "Error", err)
			return err
		}

		// Third activity: create summary from the fetched HTML
		err = workflow.ExecuteActivity(ctx, internal.GetFunctionName(activity.CreateSummary), feedItemDoc).Get(ctx, nil)
		if err != nil {
			workflow.GetLogger(ctx).Error("createSummaryActivity activity failed.", "Error", err)
			return err
		}

		workflow.GetLogger(ctx).Info("Ingest feed item workflow completed.")

		return nil
	}
}
