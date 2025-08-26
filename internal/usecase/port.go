package usecase

import (
	"worker-manager/internal/domain"
)

type WorkerUsecase interface {
	Create(saveDuration int, workerMsg string) (domain.WorkerId, error)
	Get(id domain.WorkerId) (*domain.Worker, error)
	GetData(id domain.WorkerId) (*domain.Data, error)
	Update(id domain.WorkerId, saveDuration int, workerMsg string) error
	Delete(id domain.WorkerId) error
	Count() int
}
