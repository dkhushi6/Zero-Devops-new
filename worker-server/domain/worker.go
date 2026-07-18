package domain

import "go.uber.org/zap"

type WorkerUsecase interface {
	StartWorker(baseLogger *zap.Logger) error
}