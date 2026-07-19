package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
)

func TestRequestIDMiddleware_GeneratesID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := RequestIDMiddleware(func(c *echo.Context) error {
		id := GetRequestID(c)
		if id == "" {
			t.Error("expected request ID to be set")
		}
		return nil
	})

	if err := handler(c); err != nil {
		t.Fatal(err)
	}

	respID := rec.Header().Get(RequestIDHeader)
	if respID == "" {
		t.Error("expected X-Request-Id header in response")
	}
}

func TestRequestIDMiddleware_UsesIncomingHeader(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
	req.Header.Set(RequestIDHeader, "incoming-id-123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := RequestIDMiddleware(func(c *echo.Context) error {
		id := GetRequestID(c)
		if id != "incoming-id-123" {
			t.Errorf("expected 'incoming-id-123', got '%s'", id)
		}
		return nil
	})

	if err := handler(c); err != nil {
		t.Fatal(err)
	}

	respID := rec.Header().Get(RequestIDHeader)
	if respID != "incoming-id-123" {
		t.Errorf("expected response header 'incoming-id-123', got '%s'", respID)
	}
}

func TestRequestIDMiddleware_PassesToNextHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	var called bool
	handler := RequestIDMiddleware(func(_ *echo.Context) error {
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

func TestGetRequestID_EmptyWhenNotSet(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if id := GetRequestID(c); id != "" {
		t.Errorf("expected empty string, got '%s'", id)
	}
}
