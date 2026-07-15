package domain

import (
	"context"
	"time"
)

type DeploymentStatus string

const (
	DeploymentStatusPending  DeploymentStatus = "pending"
	DeploymentStatusRunning  DeploymentStatus = "running"
	DeploymentStatusSuccess  DeploymentStatus = "success"
	DeploymentStatusFailed   DeploymentStatus = "failed"
	DeploymentStatusCanceled DeploymentStatus = "canceled"
)

type Deployment struct {
	ID        int64            `json:"id"`
	UserID    int64            `json:"user_id"`
	RepoID    int64            `json:"repo_id"`
	CloneURL  string           `json:"clone_url"`
	Status    DeploymentStatus `json:"status"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type DeploymentUsecase interface {
	CreateDeployment(ctx context.Context, userID int64, repoID int64) (*Deployment, error)
	GetDeployments(ctx context.Context, userID int64) ([]Deployment, error)
	GetDeploymentByID(ctx context.Context, userID, deploymentID int64) (*Deployment, error)
}

type DeploymentRepository interface {
	Store(ctx context.Context, d *Deployment) error
	GetByUserID(ctx context.Context, userID int64) ([]Deployment, error)
	GetByID(ctx context.Context, userID, id int64) (*Deployment, error)
	UpdateStatus(ctx context.Context, deploymentID int64, status DeploymentStatus) error
}
