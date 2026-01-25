package internal

import (
	"context"
	"errors"
	"log/slog"
)

// ParseLogLevel converts a string log level to slog.Level.
func ParseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// MultiHandler is a slog.Handler that writes to multiple Handlers.
type MultiHandler struct {
	Handlers []slog.Handler
}

func (h *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.Handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *MultiHandler) Handle(ctx context.Context, record slog.Record) error {
	var errs error
	for _, handler := range h.Handlers {
		if handler.Enabled(ctx, record.Level) {
			errs = errors.Join(errs, handler.Handle(ctx, record))
		}
	}
	return errs
}

func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h.Handlers))
	for i, handler := range h.Handlers {
		handlers[i] = handler.WithAttrs(attrs)
	}
	return &MultiHandler{Handlers: handlers}
}

func (h *MultiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(h.Handlers))
	for i, handler := range h.Handlers {
		handlers[i] = handler.WithGroup(name)
	}
	return &MultiHandler{Handlers: handlers}
}
