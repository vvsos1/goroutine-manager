package usecase

import (
	"context"
	"worker-manager/internal/domain"
)

type WorkerUsecase interface {
	Create(ctx context.Context, saveDuration int, workerMsg string) (domain.WorkerId, error)
	Get(ctx context.Context, id domain.WorkerId) (*domain.Worker, error)
	GetData(ctx context.Context, id domain.WorkerId) (*domain.Data, error)
	Update(ctx context.Context, id domain.WorkerId, saveDuration int, workerMsg string) error
	Delete(ctx context.Context, id domain.WorkerId) error
	Count(ctx context.Context) int
}
