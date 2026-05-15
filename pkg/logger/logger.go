package logger

import (
	"context"
	"log/slog"
	"os"
	"pob/pkg/tracing"
)

type ContextHandler struct {
	handler slog.Handler
}

func NewContexthandler(handler slog.Handler) *ContextHandler {
	return &ContextHandler{
		handler: handler,
	}
}

func (h *ContextHandler) Handle(ctx context.Context, record slog.Record) error {
	var attrs []slog.Attr

	trace := tracing.NewContextTrace()
	if id := trace.GetValueFromCtx(ctx); id != "" {
		attrs = append(attrs, slog.String(string(trace.Getkey()), id))
	}

	if len(attrs) > 0 {
		record.AddAttrs(attrs...)
	}
	return h.handler.Handle(ctx, record)
}

func (h *ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewContexthandler(h.handler.WithAttrs(attrs))
}

func (h *ContextHandler) WithGroup(name string) slog.Handler {
	return NewContexthandler(h.handler.WithGroup(name))
}

type contextKey struct{}

func InitLogger() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	contexthandler := NewContexthandler(handler)
	logger := slog.New(contexthandler)
	slog.SetDefault(logger)
}
