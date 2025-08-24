package usecase

import (
	"goroutine-manager/internal/domain"
)

type GoroutineUsecase interface {
	Create(saveDuration int) (domain.GoroutineId, error)
	Get(id domain.GoroutineId) (string, error)
	Update(id domain.GoroutineId, saveDuration int) error
	Delete(id domain.GoroutineId) error
	Count() int
}
