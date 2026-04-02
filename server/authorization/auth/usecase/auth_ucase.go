package usecase

import (
	"context"
	"time"
	
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"server/domain"
)

type authUsecase struct {
	userRepo domain.UserRepository
	githubRepo domain.GithubRepository
	contextTimeout time.Duration
}

func NewUserUsecase(u domain.UserRepository, g domain.GithubRepository, timeout time.Duration) domain.AuthUsecase {
	return &authUsecase{
		userRepo: u,
		githubRepo: g,
		contextTimeout: timeout,
	}
}

func (a *authUsecase) HandleGithubCallback(ctx context.Context){
	g,ctx = errgroup.WithContext(ctx)
	
	g.Go(func())
}