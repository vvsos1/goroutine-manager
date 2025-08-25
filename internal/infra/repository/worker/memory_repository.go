package worker

import (
	"fmt"
	"goroutine-manager/internal/domain"
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

func (r *MemoryRepository) Save(worker *domain.Worker) error {
	r.workers[worker.Id] = worker
	return nil
}

func (r *MemoryRepository) Get(id domain.WorkerId) (*domain.Worker, error) {
	worker, exists := r.workers[id]
	if !exists {
		return nil, fmt.Errorf("worker not exists")
	}
	return worker, nil
}

//func (r *MemoryRepository) Delete(id domain.WorkerId) error {
//	delete(r.workers, id)
//	return nil
//}

func (r *MemoryRepository) Count() int {
	return len(r.workers)
}
