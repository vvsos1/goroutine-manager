package goroutine

import (
	"fmt"
	"goroutine-manager/internal/domain"
)

var _ domain.GoroutineRepository = (*MemoryRepository)(nil)

type MemoryRepository struct {
	goroutines map[domain.GoroutineId]*domain.Goroutine
}

func NewMemoryGoroutineRepository() *MemoryRepository {
	return &MemoryRepository{
		goroutines: make(map[domain.GoroutineId]*domain.Goroutine),
	}
}

func (r *MemoryRepository) Save(goroutine *domain.Goroutine) error {
	r.goroutines[goroutine.Id] = goroutine
	return nil
}

func (r *MemoryRepository) Get(id domain.GoroutineId) (*domain.Goroutine, error) {
	goroutine, exists := r.goroutines[id]
	if !exists {
		return nil, fmt.Errorf("goroutine not exists")
	}
	return goroutine, nil
}

func (r *MemoryRepository) Delete(id domain.GoroutineId) error {
	delete(r.goroutines, id)
	return nil
}

func (r *MemoryRepository) Count() int {
	return len(r.goroutines)
}
