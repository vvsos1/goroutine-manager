package data

import (
	"fmt"
	"worker-manager/internal/domain"
)

type MemoryRepository struct {
	values map[domain.WorkerId]*domain.Data
}

var _ domain.DataRepository = (*MemoryRepository)(nil)

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		values: make(map[domain.WorkerId]*domain.Data),
	}
}

func (r *MemoryRepository) Put(key domain.WorkerId, value *domain.Data) error {
	r.values[key] = value
	return nil
}

func (r *MemoryRepository) Get(key domain.WorkerId) (*domain.Data, error) {
	if value, exists := r.values[key]; exists {
		return value, nil
	}
	return nil, fmt.Errorf("key not found: %d", key)
}

func (r *MemoryRepository) Delete(key domain.WorkerId) error {
	delete(r.values, key)
	return nil
}
