package worker

import (
	"context"
	"fmt"
	"worker-manager/internal/domain"
	"worker-manager/util/logger"
)

var _ domain.WorkerRepository = (*MemoryRepository)(nil)

type MemoryRepository struct {
	workers map[domain.WorkerId]*domain.Worker
}

func NewMemoryWorkerRepository() *MemoryRepository {
	return &MemoryRepository{
		workers: make(map[domain.WorkerId]*domain.Worker),
	}
}

func (r *MemoryRepository) Save(ctx context.Context, worker *domain.Worker) error {
	r.workers[worker.Id] = worker
	logger.Infof(ctx, "memory repository save worker: %d", worker.Id)
	return nil
}

func (r *MemoryRepository) Get(ctx context.Context, id domain.WorkerId) (*domain.Worker, error) {
	worker, exists := r.workers[id]
	if !exists {
		logger.Infof(ctx, "memory repository get worker: %d not exists", id)
		return nil, fmt.Errorf("worker not exists")
	}
	logger.Infof(ctx, "memory repository get worker: %d", id)
	return worker, nil
}

//func (r *MemoryRepository) Delete(id domain.WorkerId) error {
//	delete(r.workers, id)
//	return nil
//}

func (r *MemoryRepository) Count(ctx context.Context) int {
	logger.Infof(ctx, "memory repository count workers: %d", len(r.workers))
	return len(r.workers)
}
