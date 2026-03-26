package domain

import "context"

// Session Management Tokens
type TokenResponse struct {
	AccessToken		string 		`json:"accessToken"`
	RefreshToken	string	`json:"refreshToken"`
}


//  Oauth Methods
type AuthUsecase interface {
	HandleGithubCallback(ctx context.Context) (*TokenResponse , error)
	RefreshToken(ctx Context.context , refreshToken string) error
	Logout(t* Token, accessToken string) error
}