package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Netflix/go-env"
	"github.com/demeyerthom/feeds-aggregator/internal"
	internalactivity "github.com/demeyerthom/feeds-aggregator/internal/activity"
	internalworkflow "github.com/demeyerthom/feeds-aggregator/internal/workflow"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/opentracing"
	"go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

const serviceName = "feeds-worker"

var (
	cfg Configuration

	rdb         *redis.Client
	mongoClient *mongo.Client
	zenClient   openai.Client
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
	MongoDB struct {
		URI string `env:"MONGODB_URI,default=mongodb://localhost:27017"`
	}
	Logging struct {
		Level string `env:"LOG_LEVEL,default=info"`
	}
	Storage struct {
		HTMLDir string `env:"HTML_STORAGE_DIR,default=./data"`
	}
	TextExtractor struct {
		Limit int `env:"TEXT_LIMIT,default=400000"`
	}
	Ollama struct {
		Enabled bool   `env:"OLLAMA_ENABLED,default=false"`
		Host    string `env:"OLLAMA_HOST,default=http://localhost:11434"`
		Model   string `env:"OLLAMA_MODEL,default=gpt-oss:20b"`
	}
	OpenCode struct {
		Enabled bool   `env:"OPENCODE_ENABLED,default=false"`
		Host    string `env:"OPENCODE_HOST,default=https://opencode.ai/zen/v1"`
		Model   string `env:"OPENCODE_MODEL,default=big-pickle"`
		APIKey  string `env:"OPENCODE_API_KEY"`
	}
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
	slog.Warn("Log level set", "level", logLevel.String())
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

	// Initialize MongoDB client
	mongoCtx, mongoCancel := context.WithTimeout(ctx, 10*time.Second)
	defer mongoCancel()

	mongoClient, err = mongo.Connect(mongoCtx, options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		slog.Error("Failed to connect to MongoDB", "err", err)
		os.Exit(1)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			slog.Error("Failed to disconnect from MongoDB", "err", err)
		}
	}()

	// Ping MongoDB to verify connection
	if err := mongoClient.Ping(mongoCtx, nil); err != nil {
		slog.Error("Failed to ping MongoDB", "err", err)
		os.Exit(1)
	}
	slog.Info("Connected to MongoDB", "uri", cfg.MongoDB.URI)

	// Get collection reference and create unique index on link field
	feedItemCollection := mongoClient.Database(internal.MongoDBName).Collection(internal.MongoFeedItemCollection)
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "link", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = feedItemCollection.Indexes().CreateOne(mongoCtx, indexModel)
	if err != nil {
		slog.Error("Failed to create unique index on link field", "err", err)
		os.Exit(1)
	}

	// Validate provider configuration: exactly one must be enabled
	ollamaEnabled := cfg.Ollama.Enabled
	openCodeEnabled := cfg.OpenCode.Enabled

	if ollamaEnabled && openCodeEnabled {
		slog.Error("Both Ollama and OpenCode providers are enabled. Exactly one provider must be enabled.")
		os.Exit(1)
	}
	if !ollamaEnabled && !openCodeEnabled {
		slog.Error("No provider is enabled. Enable either Ollama (OLLAMA_ENABLED=true) or OpenCode (OPENCODE_ENABLED=true).")
		os.Exit(1)
	}

	// Initialize the configured provider and determine the model to use
	var model string
	if openCodeEnabled {
		if cfg.OpenCode.APIKey == "" {
			slog.Error("OPENCODE_API_KEY is required when OpenCode provider is enabled")
			os.Exit(1)
		}
		zenClient = openai.NewClient(
			option.WithAPIKey(cfg.OpenCode.APIKey),
			option.WithBaseURL(cfg.OpenCode.Host),
		)
		model = cfg.OpenCode.Model
		slog.Info("Initialized OpenCode client", "model", model, "host", cfg.OpenCode.Host)
	} else {
		// Ollama OpenAI compatibility: baseURL = OLLAMA_HOST + "/v1/", apiKey = "ollama" (required but ignored)
		ollamaBaseURL := cfg.Ollama.Host + "/v1/"
		zenClient = openai.NewClient(
			option.WithAPIKey("ollama"),
			option.WithBaseURL(ollamaBaseURL),
		)
		model = cfg.Ollama.Model
		slog.Info("Initialized Ollama client", "model", model, "host", cfg.Ollama.Host, "baseURL", ollamaBaseURL)
	}

	// Create interceptor
	tracingInterceptor, err := opentracing.NewInterceptor(opentracing.TracerOptions{})
	if err != nil {
		log.Fatalf("Failed creating interceptor: %v", err)
	}

	// Initialize Temporal client
	temporalClient, err := client.Dial(client.Options{
		HostPort:           cfg.Temporal.Host,
		Logger:             slog.Default(),
		ContextPropagators: []workflow.ContextPropagator{internal.NewContextPropagator()},
		Interceptors:       []interceptor.ClientInterceptor{tracingInterceptor},
	})
	if err != nil {
		slog.Error("Unable to create Temporal Client", "err", err)
		os.Exit(1)
	}
	defer temporalClient.Close()

	w := worker.New(temporalClient, internal.TaskQueueName, worker.Options{})

	// Register workflow
	w.RegisterWorkflowWithOptions(internalworkflow.IngestFeedItem(), workflow.RegisterOptions{
		Name: internal.GetFunctionName(internalworkflow.IngestFeedItem),
	})

	// Register activities with the Activities struct methods
	w.RegisterActivityWithOptions(internalactivity.AddNewFeedItem(feedItemCollection), activity.RegisterOptions{
		Name: internal.GetFunctionName(internalactivity.AddNewFeedItem),
	})
	w.RegisterActivityWithOptions(internalactivity.FetchHTML(http.DefaultClient, cfg.Storage.HTMLDir), activity.RegisterOptions{
		Name: internal.GetFunctionName(internalactivity.FetchHTML),
	})
	w.RegisterActivityWithOptions(
		internalactivity.ProcessContent(feedItemCollection, zenClient, model, cfg.Storage.HTMLDir, cfg.TextExtractor.Limit),
		activity.RegisterOptions{
			Name: internal.GetFunctionName(internalactivity.ProcessContent),
		},
	)

	slog.Info("Starting worker", "taskQueue", internal.TaskQueueName)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		slog.Error("Unable to start worker", "err", err)
		os.Exit(1)
	}
}
