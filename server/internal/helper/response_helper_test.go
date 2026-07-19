package helper

import (
	"testing"
)

func TestBuildSuccessResponse_Basic(t *testing.T) {
	data := map[string]string{"key": "value"}
	resp := BuildSuccessResponse(data, "", "req-1")

	if !resp.Success {
		t.Error("expected success to be true")
	}
	if resp.Data == nil {
		t.Fatal("expected data to be non-nil")
	}
	if resp.RequestID != "req-1" {
		t.Errorf("expected requestID 'req-1', got '%s'", resp.RequestID)
	}
	if resp.Message != "" {
		t.Errorf("expected empty message, got '%s'", resp.Message)
	}
}

func TestBuildSuccessResponse_WithMessage(t *testing.T) {
	data := "some data"
	resp := BuildSuccessResponse(data, "", "req-2", WithMessage("operation completed"))

	if !resp.Success {
		t.Error("expected success to be true")
	}
	if resp.Message != "operation completed" {
		t.Errorf("expected message 'operation completed', got '%s'", resp.Message)
	}
}

func TestBuildSuccessResponse_NilData(t *testing.T) {
	resp := BuildSuccessResponse(nil, "", "req-3")

	if resp.Data != nil {
		t.Error("expected data to be nil")
	}
}

func TestBuildSuccessResponse_MultipleOptions(t *testing.T) {
	data := 42
	resp := BuildSuccessResponse(data, "", "req-4", WithMessage("first"), WithMessage("second"))

	if resp.Message != "second" {
		t.Errorf("expected message 'second' (last wins), got '%s'", resp.Message)
	}
}
