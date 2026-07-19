package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func TestRequestLoggerMiddleware_InjectsLogger(t *testing.T) {
	viper.Set("APP_ENV", "test")
	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	base := zap.NewNop()
	handler := RequestLoggerMiddleware(base)(func(c *echo.Context) error {
		ctx := c.Request().Context()
		l := LoggerFromContext(ctx)
		if l == nil {
			t.Fatal("expected non-nil logger in context")
		}
		return nil
	})

	if err := handler(c); err != nil {
		t.Fatal(err)
	}
}

func TestLoggerFromContext_ReturnsBaseWhenNotSet(t *testing.T) {
	ctx := context.Background()
	l := LoggerFromContext(ctx)
	if l == nil {
		t.Error("expected non-nil logger, got nil")
	}
}

func TestLoggerFromContext_ReturnsLoggerFromContext(t *testing.T) {
	expected := zap.NewNop()
	ctx := context.WithValue(context.Background(), loggerCtxKey{}, expected)

	l := LoggerFromContext(ctx)
	if l != expected {
		t.Error("expected to retrieve the injected logger")
	}
}

func TestRequestLoggerMiddleware_ChainCallsNext(t *testing.T) {
	base := zap.NewNop()
	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	var called bool
	handler := RequestLoggerMiddleware(base)(func(_ *echo.Context) error {
		called = true
		return nil
	})

	if err := handler(c); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("expected next handler to be called")
	}
}
