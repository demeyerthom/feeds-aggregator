// Sources for https://watermill.io/learn/getting-started/
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/Netflix/go-env"
	"github.com/demeyerthom/feeds-aggregator/internal"
	"github.com/demeyerthom/feeds-aggregator/internal/workflow"
	"github.com/mmcdole/gofeed"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.temporal.io/sdk/client"
)

const serviceName = "feeds-ingester"

var (
	cfg Configuration

	rdb            *redis.Client
	temporalClient client.Client
	linksCounter   metric.Int64Counter
	tracer         trace.Tracer
)

type Configuration struct {
	Redis struct {
		Host     string `env:"REDIS_HOST,default=localhost:6379"`
		Password string `env:"REDIS_PASSWORD"`
		Database int    `env:"REDIS_DB,default=0"`
	}
	Otel struct {
		Host string `env:"OTEL_HOST,default=localhost:4318"`
	}
	Temporal struct {
		Host string `env:"TEMPORAL_HOST,default=localhost:7233"`
	}
	Logging struct {
		Level string `env:"LOG_LEVEL,default=info"`
	}
	SourceFile     string        `env:"SOURCE_FILE,default=/feeds.json"`
	TickerInterval time.Duration `env:"TICKER_INTERVAL,default=1m"`
}

func init() {
	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		slog.Error("Failed to unmarshal environment", "err", err)
		os.Exit(1)
	}
}

func main() {
	// Set up OTel SDK
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	otelShutdown, err := internal.SetupOTelSDK(ctx, serviceName, cfg.Otel.Host)
	if err != nil {
		slog.Error("Failed to setup OTel SDK", "err", err)
		os.Exit(1)
	}
	defer func() {
		if err := otelShutdown(context.Background()); err != nil {
			slog.Error("Failed to shutdown OTel SDK", "err", err)
		}
	}()

	logLevel := internal.ParseLogLevel(cfg.Logging.Level)
	slog.Debug("Log level set", "level", logLevel.String())
	slog.SetLogLoggerLevel(logLevel)
	logger := slog.New(&internal.MultiHandler{
		Handlers: []slog.Handler{
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: logLevel,
			}),
			otelslog.NewHandler(serviceName),
		},
	})
	slog.SetDefault(logger)

	// Initialize metrics
	meter := otel.Meter(serviceName)
	linksCounter, err = meter.Int64Counter(
		"feeds.new_links",
		metric.WithDescription("Number of new feed links discovered and stored in Redis"),
		metric.WithUnit("{link}"),
	)
	if err != nil {
		slog.Error("Failed to create metric counter", "err", err)
		os.Exit(1)
	}

	// Initialize tracer
	tracer = otel.Tracer(serviceName)

	// Initialize Redis client
	rdb = redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Host,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.Database,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			slog.Debug("New Redis connection established")
			return nil
		},
	})
	defer func() {
		if err := rdb.Close(); err != nil {
			slog.Error("Failed to close Redis client", "err", err)
		}
	}()
	go func() {
		for {
			stats := rdb.PoolStats()
			slog.Debug("Redis pool stats",
				"hits", stats.Hits,
				"misses", stats.Misses,
				"timeouts", stats.Timeouts,
				"totalConns", stats.TotalConns,
				"idleConns", stats.IdleConns,
			)
			time.Sleep(30 * time.Second)
		}
	}()

	// Initialize Temporal client
	temporalClient, err = client.Dial(client.Options{
		HostPort: cfg.Temporal.Host,
		Logger:   slog.Default(),
	})
	if err != nil {
		slog.Error("Unable to create Temporal Client", "err", err)
		os.Exit(1)
	}
	defer temporalClient.Close()

	// Load feed list from source file
	feedList, err := loadFeedList(cfg.SourceFile)
	if err != nil {
		slog.Error("Failed to load feed list", "err", err)
		os.Exit(1)
	}
	slog.Info("Loaded feed list", "count", len(feedList))

	processAllFeeds(ctx, feedList)

	// Set up ticker to process feeds every 5 minutes
	ticker := time.NewTicker(cfg.TickerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Shutting down feed worker")
			return
		case <-ticker.C:
			processAllFeeds(ctx, feedList)
		}
	}
}

func processAllFeeds(ctx context.Context, feedList internal.FeedList) {
	slog.Info("Starting feed processing cycle", "count", len(feedList))
	for _, feed := range feedList {
		if err := processFeedActivity(ctx, feed); err != nil {
			slog.Error("Failed to process feed", "feed", feed.Title, "err", err)
		}
	}
	slog.Info("Completed feed processing cycle")
}

func loadFeedList(path string) (internal.FeedList, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var feedList internal.FeedList
	if err := json.Unmarshal(b, &feedList); err != nil {
		return nil, err
	}

	return feedList, nil
}

func processFeedActivity(ctx context.Context, f internal.Feed) error {
	// Start a span for the entire activity
	ctx, span := tracer.Start(ctx, "processFeed",
		trace.WithAttributes(
			attribute.String("feed.title", f.Title),
			attribute.String("feed.url", f.XMLURL),
		),
	)
	defer span.End()

	fp := gofeed.NewParser()

	// Trace the feed parsing (HTTP call)
	parseCtx, parseSpan := tracer.Start(ctx, "parseFeedURL",
		trace.WithAttributes(attribute.String("url", f.XMLURL)),
	)
	feed, err := fp.ParseURL(f.XMLURL)
	if err != nil {
		parseSpan.RecordError(err)
		parseSpan.SetStatus(codes.Error, "failed to parse feed")
		parseSpan.End()
		slog.Error("Unable to parse feed URL", "URL", f.XMLURL, "err", err)
		return nil
	}
	parseSpan.SetAttributes(attribute.Int("items.count", len(feed.Items)))
	parseSpan.End()

	// Process each item
	for _, item := range feed.Items {
		processItem(parseCtx, feed.Title, item)
	}

	return nil
}

func processItem(ctx context.Context, feedTitle string, item *gofeed.Item) {
	ctx, span := tracer.Start(ctx, "processItem",
		trace.WithAttributes(
			attribute.String("item.link", item.Link),
			attribute.String("item.title", item.Title),
		),
	)
	defer span.End()

	_, err := rdb.Get(ctx, item.Link).Result()
	if errors.Is(err, redis.Nil) {
		// New item - store in Redis
		rdb.Set(ctx, item.Link, "1", 24*7*time.Hour)
		linksCounter.Add(ctx, 1, metric.WithAttributeSet(
			attribute.NewSet(attribute.String("feed.title", feedTitle))))
		span.SetAttributes(attribute.Bool("item.new", true))

		// Start Temporal workflow for this feed item
		feedItem := internal.FeedItem{
			Link:  item.Link,
			Title: item.Title,
		}
		workflowOptions := client.StartWorkflowOptions{
			ID:        fmt.Sprintf("ingest-feed-item-%s", item.Link),
			TaskQueue: internal.TaskQueueName,
		}
		we, err := temporalClient.ExecuteWorkflow(ctx, workflowOptions, internal.GetFunctionName(workflow.IngestFeedItem), feedItem)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "temporal workflow error")
			slog.Error("Failed to start Temporal workflow", "err", err)
			return
		}

		slog.Info("Started workflow for feed item", "workflowID", we.GetID(), "runID", we.GetRunID(), "link", item.Link)
	} else if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "redis error")
		slog.Error("Failed to get item from Redis", "err", err.Error())
	} else {
		span.SetAttributes(attribute.Bool("item.new", false))
		slog.Debug(fmt.Sprintf("Item already processed: %s", item.Link))
	}
}
