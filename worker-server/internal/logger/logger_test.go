package logger

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestNew_Production(t *testing.T) {
	l, err := New("production")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
	if l.Core().Enabled(zapcore.DebugLevel) {
		t.Error("expected debug level to be disabled in production")
	}
}

func TestNew_Development(t *testing.T) {
	l, err := New("development")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
	if !l.Core().Enabled(zapcore.DebugLevel) {
		t.Error("expected debug level to be enabled in development")
	}
}

func TestNew_DefaultToDevelopment(t *testing.T) {
	l, err := New("staging")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
	if !l.Core().Enabled(zapcore.DebugLevel) {
		t.Error("expected debug level to be enabled by default")
	}
}
