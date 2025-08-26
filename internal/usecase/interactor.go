package usecase

import (
	"context"
	"fmt"
	"worker-manager/internal/domain"
	"worker-manager/util/logger"
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

func (gm *WorkerInteractor) Create(ctx context.Context, saveDuration int, workerMsg string) (domain.WorkerId, error) {
	worker := domain.NewWorker(saveDuration, workerMsg, gm.dataRepository)
	err := gm.workerRepository.Save(ctx, worker)
	if err != nil {
		logger.Errorf(ctx, "failed to create worker: %v", err)
		return -1, fmt.Errorf("failed to create worker: %v", err)
	}
	worker.StartInGoroutine()
	logger.Infof(ctx, "created worker %v", worker.Id)
	return worker.Id, nil
}

func (gm *WorkerInteractor) Get(ctx context.Context, id domain.WorkerId) (*domain.Worker, error) {
	worker, err := gm.workerRepository.Get(ctx, id)
	if err != nil {
		logger.Errorf(ctx, "failed to get worker: %v", err)
		return nil, fmt.Errorf("worker with id %d not found", id)
	}
	logger.Infof(ctx, "get worker %v", id)
	return worker, nil
}

func (gm *WorkerInteractor) GetData(ctx context.Context, id domain.WorkerId) (*domain.Data, error) {
	worker, err := gm.workerRepository.Get(ctx, id)
	if err != nil {
		logger.Errorf(ctx, "failed to get data for worker: %v", err)
		return nil, fmt.Errorf("worker with id %d not found", id)
	}
	value := worker.Read()
	logger.Infof(ctx, "get data for worker %v", id)
	return value, nil
}

func (gm *WorkerInteractor) Update(ctx context.Context, id domain.WorkerId, saveDuration int, workerMsg string) error {
	worker, err := gm.workerRepository.Get(ctx, id)

	if err != nil {
		logger.Errorf(ctx, "failed to update worker: %v; not found", err)
		return fmt.Errorf("worker with id %d not found", id)
	}
	worker.Update(saveDuration, workerMsg)

	err = gm.workerRepository.Save(ctx, worker)
	if err != nil {
		logger.Errorf(ctx, "failed to update worker: %v; not save", err)
		return fmt.Errorf("worker with id %d was not updated", id)
	}
	logger.Infof(ctx, "updated worker %v", worker.Id)
	return nil
}

func (gm *WorkerInteractor) Delete(ctx context.Context, id domain.WorkerId) error {
	worker, err := gm.workerRepository.Get(ctx, id)
	if err != nil {
		logger.Errorf(ctx, "failed to delete worker: %v; not found", err)
		return fmt.Errorf("worker with id %d not found", id)
	}
	worker.Delete()

	//err = gm.workerRepository.Delete(id)
	//if err != nil {
	//	return fmt.Errorf("worker with id %d was not deleted", id)
	//}

	logger.Infof(ctx, "deleted worker %v", worker.Id)
	return nil
}

func (gm *WorkerInteractor) Count(ctx context.Context) int {
	count := gm.workerRepository.Count(ctx)
	logger.Infof(ctx, "current worker count: %d", count)
	return count
}
