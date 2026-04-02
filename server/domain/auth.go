package domain

import "context"

// Session Management Tokens
type TokenResponse struct {
	AccessToken		string 		`json:"accessToken"`
	RefreshToken	string	`json:"refreshToken"`
}

type LoginApps struct {
	Apps map[string]string `json:"apps"`
}

//  Oauth Methods
type AuthUsecase interface {
	Signup(ctx context.Context , provider string) error
	Login(ctx context.Context, provider string) error
	RefreshToken(ctx context.Context, refreshToken string) error
	Logout(ctx context.Context, accessToken string) error
}