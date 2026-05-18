package usecase

import (
	"context"
	"net/http"
	"time"
	"Zero_Devops/server/domain"
)

type githubAppUsecase struct{
	// Here gitRepo does not identify github repositories it means the github functions 
	githubRepo domain.GithubRepository
}

func NewGithubAppUsecase() *githubAppUsecase {
	return &githubAppUsecase{}
}

func (g *githubAppUsecase) InstallGithubApp(ctx context.Context) error {

	return nil
}

func (g *githubAppUsecase) GetGithubAppInstallation(ctx context.Context,user_id int64) error {
	return nil
}

func (g *githubAppUsecase) DeleteGithubApp(ctx context.Context) error {
	return nil
}
