package http

import (
	"Zero_Devops/server/internal/domain"
	"Zero_Devops/server/internal/helper"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/spf13/viper"
)

type mockAuthUsecase struct {
	tokens         *domain.TokenResponse
	userResp       *domain.UserResponse
	err            error
	githubOAuthURL string
}

func (m *mockAuthUsecase) HandleOAuthCallback(_ context.Context, _, _ string) (*domain.TokenResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.tokens, nil
}

func (m *mockAuthUsecase) RefreshToken(_ context.Context, _ string) (*domain.TokenResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.tokens, nil
}

func (m *mockAuthUsecase) GetCurrentUser(_ context.Context, _ string) (domain.UserResponse, error) {
	if m.err != nil {
		return domain.UserResponse{}, m.err
	}
	return *m.userResp, nil
}

func (m *mockAuthUsecase) Logout(_ context.Context, _ string) error {
	return m.err
}

func (m *mockAuthUsecase) GithubOauthURL(_ context.Context, _ string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.githubOAuthURL, nil
}

func setTestConfig() {
	viper.Set("JWT_SECRET", "test-secret-key")
	viper.Set("IS_PRODUCTION_ENV", false)
	viper.Set("ACCESS_TOKEN_EXPIRY", "1")
	viper.Set("REFRESH_TOKEN_EXPIRY", "720")
}

func TestLogin_DefaultReturnTo(t *testing.T) {
	setTestConfig()

	handler := &AuthHandler{
		AUsecase: &mockAuthUsecase{
			githubOAuthURL: "https://github.com/login/oauth/authorize?state=test-state",
		},
	}

	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/auth/github/login", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Login(c)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if rec.Code != http.StatusTemporaryRedirect {
		t.Errorf("expected status %d, got %d", http.StatusTemporaryRedirect, rec.Code)
	}

	location := rec.Header().Get("Location")
	if location != "https://github.com/login/oauth/authorize?state=test-state" {
		t.Errorf("expected redirect to GitHub OAuth URL, got %s", location)
	}
}

func TestLogin_Success(t *testing.T) {
	setTestConfig()

	handler := &AuthHandler{
		AUsecase: &mockAuthUsecase{
			githubOAuthURL: "https://github.com/login/oauth/authorize?state=test-state",
		},
	}

	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/auth/github/login?return_to=/dashboard", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Login(c)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if rec.Code != http.StatusTemporaryRedirect {
		t.Errorf("expected status %d, got %d", http.StatusTemporaryRedirect, rec.Code)
	}

	location := rec.Header().Get("Location")
	if location != "https://github.com/login/oauth/authorize?state=test-state" {
		t.Errorf("expected redirect to GitHub OAuth URL, got %s", location)
	}

	cookies := rec.Result().Cookies()
	var hasStateCookie bool
	for _, cookie := range cookies {
		if cookie.Name == "gh_oauth_state" {
			hasStateCookie = true
			break
		}
	}
	if !hasStateCookie {
		t.Error("expected gh_oauth_state cookie to be set")
	}
}

func TestLogin_UsecaseError(t *testing.T) {
	setTestConfig()

	handler := &AuthHandler{
		AUsecase: &mockAuthUsecase{
			err: domain.ErrInternalServerError,
		},
	}

	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/auth/github/login", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Login(c)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestRefresh_MissingCookie(t *testing.T) {
	setTestConfig()

	handler := &AuthHandler{
		AUsecase: &mockAuthUsecase{},
	}

	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Refresh(c)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestRefresh_Success(t *testing.T) {
	setTestConfig()

	handler := &AuthHandler{
		AUsecase: &mockAuthUsecase{
			tokens: &domain.TokenResponse{
				AccessToken:  "new-access-token",
				RefreshToken: "new-refresh-token",
			},
		},
	}

	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", http.NoBody)
	//nolint:gosec
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "test-refresh-token"})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Refresh(c)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestRefresh_UsecaseError(t *testing.T) {
	setTestConfig()

	handler := &AuthHandler{
		AUsecase: &mockAuthUsecase{
			err: domain.ErrInvalidToken,
		},
	}

	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", http.NoBody)
	//nolint:gosec
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "invalid-token"})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Refresh(c)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestLogout_MissingCookie(t *testing.T) {
	setTestConfig()

	handler := &AuthHandler{
		AUsecase: &mockAuthUsecase{},
	}

	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Logout(c)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestLogout_Success(t *testing.T) {
	setTestConfig()

	handler := &AuthHandler{
		AUsecase: &mockAuthUsecase{},
	}

	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", http.NoBody)
	//nolint:gosec
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "test-access-token"})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Logout(c)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var accessCookie *http.Cookie
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == "access_token" {
			accessCookie = cookie
			break
		}
	}
	if accessCookie == nil || accessCookie.MaxAge != -1 {
		t.Error("expected access_token cookie to be cleared")
	}
}

func TestLogout_UsecaseError(t *testing.T) {
	setTestConfig()

	handler := &AuthHandler{
		AUsecase: &mockAuthUsecase{
			err: domain.ErrInvalidToken,
		},
	}

	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", http.NoBody)
	//nolint:gosec
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "test-access-token"})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Logout(c)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestGetUser_MissingCookie(t *testing.T) {
	setTestConfig()

	handler := &AuthHandler{
		AUsecase: &mockAuthUsecase{},
	}

	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetUser(c)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestGetUser_Success(t *testing.T) {
	setTestConfig()

	handler := &AuthHandler{
		AUsecase: &mockAuthUsecase{
			userResp: &domain.UserResponse{
				ID:        "1",
				Provider:  "github",
				Username:  "testuser",
				Email:     "test@example.com",
				AvatarURL: "https://example.com/avatar.png",
			},
		},
	}

	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
	//nolint:gosec
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "test-access-token"})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetUser(c)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "testuser") {
		t.Error("expected response to contain username")
	}
}

func TestGetUser_UsecaseError(t *testing.T) {
	setTestConfig()

	handler := &AuthHandler{
		AUsecase: &mockAuthUsecase{
			err: domain.ErrNotFound,
		},
	}

	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
	//nolint:gosec
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "test-access-token"})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetUser(c)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func validStateCookie(returnTo string) *http.Cookie {
	payload := oauthStatePayload{State: "test-state", ReturnTo: returnTo}
	b, _ := json.Marshal(payload)
	//nolint:gosec
	return &http.Cookie{Name: "gh_oauth_state", Value: base64.URLEncoding.EncodeToString(b), HttpOnly: true, SameSite: http.SameSiteLaxMode}
}

func TestLoginCallback_MissingStateCookie(t *testing.T) {
	setTestConfig()

	handler := &AuthHandler{
		AUsecase: &mockAuthUsecase{},
	}

	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/auth/github/login/callback?state=test-state&code=test-code", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.LoginCallback(c)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestLoginCallback_InvalidBase64Cookie(t *testing.T) {
	setTestConfig()

	handler := &AuthHandler{
		AUsecase: &mockAuthUsecase{},
	}

	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/auth/github/login/callback?state=test-state&code=test-code", http.NoBody)
	//nolint:gosec
	req.AddCookie(&http.Cookie{Name: "gh_oauth_state", Value: "!!!invalid-base64!!!"})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.LoginCallback(c)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestLoginCallback_InvalidJSONCookie(t *testing.T) {
	setTestConfig()

	handler := &AuthHandler{
		AUsecase: &mockAuthUsecase{},
	}

	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/auth/github/login/callback?state=test-state&code=test-code", http.NoBody)
	//nolint:gosec
	req.AddCookie(&http.Cookie{Name: "gh_oauth_state", Value: base64.URLEncoding.EncodeToString([]byte("not-json"))})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.LoginCallback(c)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestLoginCallback_StateMismatch(t *testing.T) {
	setTestConfig()

	handler := &AuthHandler{
		AUsecase: &mockAuthUsecase{},
	}

	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/auth/github/login/callback?state=wrong-state&code=test-code", http.NoBody)
	req.AddCookie(validStateCookie("/"))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.LoginCallback(c)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, rec.Code)
	}
}

func TestLoginCallback_MissingCode(t *testing.T) {
	setTestConfig()

	handler := &AuthHandler{
		AUsecase: &mockAuthUsecase{},
	}

	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/auth/github/login/callback?state=test-state", http.NoBody)
	req.AddCookie(validStateCookie("/"))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.LoginCallback(c)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestLoginCallback_Success(t *testing.T) {
	setTestConfig()

	handler := &AuthHandler{
		AUsecase: &mockAuthUsecase{
			tokens: &domain.TokenResponse{
				AccessToken:  "test-access-token",
				RefreshToken: "test-refresh-token",
			},
		},
	}

	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/auth/github/login/callback?state=test-state&code=test-code", http.NoBody)
	req.AddCookie(validStateCookie("/dashboard"))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.LoginCallback(c)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if rec.Code != http.StatusTemporaryRedirect {
		t.Errorf("expected status %d, got %d", http.StatusTemporaryRedirect, rec.Code)
	}

	location := rec.Header().Get("Location")
	if location != "/dashboard" {
		t.Errorf("expected redirect to /dashboard, got %s", location)
	}

	if len(rec.Result().Cookies()) != 2 {
		t.Errorf("expected 2 cookies, got %d", len(rec.Result().Cookies()))
	}
}

func TestLoginCallback_UsecaseError(t *testing.T) {
	setTestConfig()

	handler := &AuthHandler{
		AUsecase: &mockAuthUsecase{
			err: domain.ErrInternalServerError,
		},
	}

	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/auth/github/login/callback?state=test-state&code=test-code", http.NoBody)
	req.AddCookie(validStateCookie("/"))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.LoginCallback(c)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestGetStatusCode(t *testing.T) {
	tests := []struct {
		err    error
		expect int
	}{
		{nil, http.StatusOK},
		{domain.ErrInternalServerError, http.StatusInternalServerError},
		{domain.ErrNotFound, http.StatusNotFound},
		{domain.ErrConflict, http.StatusConflict},
		{domain.ErrInvalidToken, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		result := helper.GetStatusCode(tt.err)
		if result != tt.expect {
			t.Errorf("expected %d, got %d for error %v", tt.expect, result, tt.err)
		}
	}
}
