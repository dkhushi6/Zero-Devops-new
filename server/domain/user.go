package domain

import (
	"context"
	"time"
)


// Represents the User data struct
type User struct {
	ID			int64		`json:"id"`
	GithubID	int64	`json:"githubId"`
	Username	string		`json:"username"`
	Email 		string	`json:"email"`
	AvatarURL 	string	`json:"avatarURL"`
	CreatedAt 	time.Time	`json:"createdAt"`
}


// Represent the user's usecases
type UserUsecase interface {
	GetByID(ctx context.Context , id int64) (User , error)
	Register(ctx  context.Context , u* User) error
}

type UserRepository interface {
	GetByID(ctx context.Context , id int64) (User, error)
	GetByUsername(ctx context.Context , username string) (User , error)
	Store(ctx context.Context,u *User) error
}