package usecase

import (
	"fmt"
	"goroutine-manager/internal/domain"
)

type WorkerInteractor struct {
	workerRepository domain.WorkerRepository
	dataRepository   domain.DataRepository
}

var _ WorkerUsecase = (*WorkerInteractor)(nil)

func NewWorkerInteractor(workerRepository domain.WorkerRepository, dataRepository domain.DataRepository) *WorkerInteractor {
	return &WorkerInteractor{
		workerRepository: workerRepository,
		dataRepository:   dataRepository,
	}
}

func (gm *WorkerInteractor) Create(saveDuration int, workerMsg string) (domain.WorkerId, error) {
	worker := domain.NewWorker(saveDuration, workerMsg, gm.dataRepository)
	err := gm.workerRepository.Save(worker)
	if err != nil {
		return -1, fmt.Errorf("failed to create worker: %v", err)
	}

	worker.StartInGoroutine()

	return worker.Id, nil
}

func (gm *WorkerInteractor) Get(id domain.WorkerId) (*domain.Worker, error) {
	worker, err := gm.workerRepository.Get(id)
	if err != nil {
		return nil, fmt.Errorf("worker with id %d not found", id)
	}
	return worker, nil
}

func (gm *WorkerInteractor) GetData(id domain.WorkerId) (*domain.Data, error) {
	worker, err := gm.workerRepository.Get(id)
	if err != nil {
		return nil, fmt.Errorf("worker with id %d not found", id)
	}
	value := worker.Read()
	return value, nil
}

func (gm *WorkerInteractor) Update(id domain.WorkerId, saveDuration int, workerMsg string) error {
	worker, err := gm.workerRepository.Get(id)

	if err != nil {
		return fmt.Errorf("worker with id %d not found", id)
	}
	worker.Update(saveDuration, workerMsg)

	err = gm.workerRepository.Save(worker)
	if err != nil {
		return fmt.Errorf("worker with id %d was not updated", id)
	}

	return nil
}

func (gm *WorkerInteractor) Delete(id domain.WorkerId) error {
	worker, err := gm.workerRepository.Get(id)
	if err != nil {
		return fmt.Errorf("worker with id %d not found", id)
	}
	worker.Delete()

	//err = gm.workerRepository.Delete(id)
	//if err != nil {
	//	return fmt.Errorf("worker with id %d was not deleted", id)
	//}

	return nil
}

func (gm *WorkerInteractor) Count() int {
	return gm.workerRepository.Count()
}
