package activity

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/demeyerthom/feeds-aggregator/internal"
	"go.temporal.io/sdk/activity"
)

// FetchHTML fetches the HTML page from the feed item's link and stores it on disk.
//
// @param ctx - Context for the activity
// @param feedItemDoc - The FeedItemDocument containing the link and ID
// @return error - Returns an error if fetching or storing fails
// @author GitHub Copilot
func FetchHTML(httpClient *http.Client, dataDir string) func(ctx context.Context, feedItemDoc internal.FeedItemDocument) error {
	return func(ctx context.Context, feedItemDoc internal.FeedItemDocument) error {
		logger := activity.GetLogger(ctx)

		if err := os.MkdirAll(dataDir, 0755); err != nil {
			logger.Error("Failed to create HTML storage directory", "err", err, "dir", dataDir)
			return err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedItemDoc.Link, nil)
		if err != nil {
			logger.Error("Failed to create HTTP request", "err", err, "link", feedItemDoc.Link)
			return err
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; FeedsAggregator/1.0)")

		resp, err := httpClient.Do(req)
		if err != nil {
			logger.Error("Failed to fetch HTML page", "err", err, "link", feedItemDoc.Link)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			logger.Error("Received non-OK HTTP status", "status", resp.StatusCode, "link", feedItemDoc.Link)
			return errors.New("received non-OK HTTP status: " + resp.Status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error("Failed to read response body", "err", err, "link", feedItemDoc.Link)
			return err
		}

		filename := filepath.Join(dataDir, feedItemDoc.ID.Hex()+".html")

		if err := os.WriteFile(filename, body, 0644); err != nil {
			logger.Error("Failed to write HTML to disk", "err", err, "filename", filename)
			return err
		}

		logger.Info("Successfully fetched and stored HTML", "id", feedItemDoc.ID.Hex(), "link", feedItemDoc.Link, "filename", filename, "size", len(body))
		return nil
	}

}
