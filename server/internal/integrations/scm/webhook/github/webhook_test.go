package github

import (
	"Zero_Devops/server/internal/domain"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func withSecret(secret string) Option {
	return func(h *Webhook) error {
		h.secret = secret
		return nil
	}
}

func TestNew_NoOptions(t *testing.T) {
	hook, err := New()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if hook == nil {
		t.Fatal("expected non-nil webhook")
	}
}

func TestNew_WithSecretOption(t *testing.T) {
	hook, err := New(withSecret("mysecret"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if hook.secret != "mysecret" {
		t.Errorf("expected secret 'mysecret', got '%s'", hook.secret)
	}
}

func TestParse_NoEventsSpecified(t *testing.T) {
	hook, _ := New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", http.NoBody)

	_, err := hook.Parse(req)
	if !errors.Is(err, domain.ErrEventNotSpecifiedToParse) {
		t.Errorf("expected ErrEventNotSpecifiedToParse, got %v", err)
	}
}

func TestParse_InvalidHTTPMethod(t *testing.T) {
	hook, _ := New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)

	_, err := hook.Parse(req, domain.PushEventP)
	if !errors.Is(err, domain.ErrInvalidHTTPMethod) {
		t.Errorf("expected ErrInvalidHTTPMethod, got %v", err)
	}
}

func TestParse_MissingGithubEventHeader(t *testing.T) {
	hook, _ := New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", http.NoBody)

	_, err := hook.Parse(req, domain.PushEventP)
	if !errors.Is(err, domain.ErrMissingGithubEventHeader) {
		t.Errorf("expected ErrMissingGithubEventHeader, got %v", err)
	}
}

func TestParse_EventNotFound(t *testing.T) {
	hook, _ := New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", http.NoBody)
	req.Header.Set("X-GitHub-Event", "unknown_event")

	_, err := hook.Parse(req, domain.PushEventP)
	if !errors.Is(err, domain.ErrEventNotFound) {
		t.Errorf("expected ErrEventNotFound, got %v", err)
	}
}

func TestParse_EmptyBody(t *testing.T) {
	hook, _ := New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", http.NoBody)
	req.Header.Set("X-GitHub-Event", "push")

	_, err := hook.Parse(req, domain.PushEventP)
	if !errors.Is(err, domain.ErrParsingPayload) {
		t.Errorf("expected ErrParsingPayload, got %v", err)
	}
}

func TestParse_MissingSignatureWhenSecretIsSet(t *testing.T) {
	hook, _ := New(withSecret("secret"))

	body := `{"ref":"refs/heads/main"}`
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("Content-Type", "application/json")

	_, err := hook.Parse(req, domain.PushEventP)
	if !errors.Is(err, domain.ErrMissingHubSignatureHeader) {
		t.Errorf("expected ErrMissingHubSignatureHeader, got %v", err)
	}
}

func TestParse_HMACVerificationFailed(t *testing.T) {
	hook, _ := New(withSecret("secret"))

	body := `{"ref":"refs/heads/main"}`
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-Hub-Signature-256", "sha256=invalid_signature")

	_, err := hook.Parse(req, domain.PushEventP)
	if !errors.Is(err, domain.ErrHMACVerificationFailed) {
		t.Errorf("expected ErrHMACVerificationFailed, got %v", err)
	}
}

func signPayload(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func TestParse_PushEvent_Success(t *testing.T) {
	hook, _ := New(withSecret("secret"))

	jsonPayload := `{
		"ref": "refs/heads/main",
		"head_commit": {
			"id": "abc123",
			"message": "test commit",
			"timestamp": "2026-01-01T00:00:00Z",
			"author": {"name": "t", "email": "t@t.com", "username": "t"},
			"committer": {"name": "t", "email": "t@t.com", "username": "t"}
		},
		"repository": {
			"id": 1,
			"name": "test-repo",
			"full_name": "user/test-repo",
			"clone_url": "https://github.com/user/test-repo.git"
		}
	}`
	body := []byte(jsonPayload)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", strings.NewReader(jsonPayload))
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-Hub-Signature-256", signPayload(body, "secret"))

	result, err := hook.Parse(req, domain.PushEventP)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, ok := result.(domain.PushPayload)
	if !ok {
		t.Fatalf("expected PushPayload, got %T", result)
	}
	if got.Ref != "refs/heads/main" {
		t.Errorf("expected ref 'refs/heads/main', got '%s'", got.Ref)
	}
	if got.Repository.CloneURL != "https://github.com/user/test-repo.git" {
		t.Errorf("expected clone_url 'https://github.com/user/test-repo.git', got '%s'", got.Repository.CloneURL)
	}
}

func TestParse_InstallationEvent_Success(t *testing.T) {
	hook, _ := New()

	jsonPayload := `{"action":"created","installation":{"id":42}}`

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", strings.NewReader(jsonPayload))
	req.Header.Set("X-GitHub-Event", "installation")

	result, err := hook.Parse(req, domain.InstallationEvent)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, ok := result.(domain.InstallationPayload)
	if !ok {
		t.Fatalf("expected InstallationPayload, got %T", result)
	}
	if got.Action != "created" {
		t.Errorf("expected action 'created', got '%s'", got.Action)
	}
	if got.Installation.ID != 42 {
		t.Errorf("expected installation ID 42, got %d", got.Installation.ID)
	}
}

func TestParse_InstallationRepositoriesEvent_Success(t *testing.T) {
	hook, _ := New()

	jsonPayload := `{
		"action": "added",
		"installation": {"id": 42},
		"repositories_added": [{"id": 100, "name": "repo1", "full_name": "user/repo1"}]
	}`

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", strings.NewReader(jsonPayload))
	req.Header.Set("X-GitHub-Event", "installation_repositories")

	result, err := hook.Parse(req, domain.InstallationRepositoriesEvent)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, ok := result.(domain.InstallationRepositoriesPayload)
	if !ok {
		t.Fatalf("expected InstallationRepositoriesPayload, got %T", result)
	}
	if got.Action != "added" {
		t.Errorf("expected action 'added', got '%s'", got.Action)
	}
}

func TestParse_UnknownEvent(t *testing.T) {
	hook, _ := New()

	jsonPayload := `{"action":"unknown"}`
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", strings.NewReader(jsonPayload))
	req.Header.Set("X-GitHub-Event", "some_new_event")

	_, err := hook.Parse(req, domain.Event("some_new_event"))
	if err == nil {
		t.Fatal("expected error for unknown event")
	}
	if !strings.Contains(err.Error(), "unknown event") {
		t.Errorf("expected 'unknown event' in error, got '%s'", err.Error())
	}
}

func TestParse_ValidatesSignatureCorrectly(t *testing.T) {
	hook, _ := New(withSecret("my_secret_key"))

	body := []byte(`{"ref":"refs/heads/main"}`)
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", strings.NewReader(string(body)))
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-Hub-Signature-256", signPayload(body, "my_secret_key"))

	_, err := hook.Parse(req, domain.PushEventP)
	if err != nil {
		t.Fatalf("expected no error with correct signature, got %v", err)
	}
}

func TestParse_RejectsWrongSecret(t *testing.T) {
	hook, _ := New(withSecret("correct_secret"))

	body := []byte(`{"ref":"refs/heads/main"}`)
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", strings.NewReader(string(body)))
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-Hub-Signature-256", signPayload(body, "wrong_secret"))

	_, err := hook.Parse(req, domain.PushEventP)
	if !errors.Is(err, domain.ErrHMACVerificationFailed) {
		t.Errorf("expected ErrHMACVerificationFailed, got %v", err)
	}
}
