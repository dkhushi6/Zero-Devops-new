package middleware

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

func TestWithLoggerStoresInContext(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()

	ctx = WithLogger(ctx, logger)
	got := LoggerFromContext(ctx)
	if got != logger {
		t.Error("WithLogger did not store the logger in context")
	}
}

func TestLoggerFromContext_ReturnsBaseWhenNotSet(t *testing.T) {
	ctx := context.Background()
	l := LoggerFromContext(ctx)
	if l == nil {
		t.Error("expected non-nil logger, got nil")
	}
}

func TestLoggerFromContext_ReturnsStoredLogger(t *testing.T) {
	expected := zap.NewNop()
	ctx := context.WithValue(context.Background(), loggerCtxKey{}, expected)

	l := LoggerFromContext(ctx)
	if l != expected {
		t.Error("expected to retrieve the stored logger")
	}
}
