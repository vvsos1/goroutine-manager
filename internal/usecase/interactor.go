package usecase

import (
	"fmt"
	"goroutine-manager/internal/domain"
)

type GoroutineInteractor struct {
	goroutineRepository domain.GoroutineRepository
	dataRepository      domain.KeyValueRepository
}

var _ GoroutineUsecase = (*GoroutineInteractor)(nil)

func NewGoroutineInteractor(goroutineRepository domain.GoroutineRepository, repository domain.KeyValueRepository) *GoroutineInteractor {
	return &GoroutineInteractor{
		goroutineRepository: goroutineRepository,
		dataRepository:      repository,
	}
}

func (gm *GoroutineInteractor) Create(saveDuration int) (domain.GoroutineId, error) {
	goroutine := domain.NewGoroutine(saveDuration, gm.dataRepository)
	err := gm.goroutineRepository.Save(goroutine)
	if err != nil {
		return -1, fmt.Errorf("failed to create goroutine: %v", err)
	}

	goroutine.StartInGoroutine()

	return goroutine.Id, nil
}

func (gm *GoroutineInteractor) Get(id domain.GoroutineId) (string, error) {
	goroutine, err := gm.goroutineRepository.Get(id)
	if err != nil {
		return "", fmt.Errorf("goroutine with id %d not found", id)
	}
	value := goroutine.Read()
	return value, nil
}

func (gm *GoroutineInteractor) Update(id domain.GoroutineId, saveDuration int) error {
	goroutine, err := gm.goroutineRepository.Get(id)

	if err != nil {
		return fmt.Errorf("goroutine with id %d not found", id)
	}
	goroutine.Update(saveDuration)

	err = gm.goroutineRepository.Save(goroutine)
	if err != nil {
		return fmt.Errorf("goroutine with id %d was not updated", id)
	}

	return nil
}

func (gm *GoroutineInteractor) Delete(id domain.GoroutineId) error {
	goroutine, err := gm.goroutineRepository.Get(id)
	if err != nil {
		return fmt.Errorf("goroutine with id %d not found", id)
	}
	goroutine.Delete()

	err = gm.goroutineRepository.Delete(id)
	if err != nil {
		return fmt.Errorf("goroutine with id %d was not deleted", id)
	}

	return nil
}

func (gm *GoroutineInteractor) Count() int {
	return gm.goroutineRepository.Count()
}
