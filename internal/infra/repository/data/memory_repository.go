package data

import (
	"fmt"
	"goroutine-manager/internal/domain"
)

type MemoryRepository struct {
	values map[domain.GoroutineId]string
}

var _ domain.KeyValueRepository = (*MemoryRepository)(nil)

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		values: make(map[domain.GoroutineId]string),
	}
}

func (r *MemoryRepository) Put(key domain.GoroutineId, value string) error {
	r.values[key] = value
	return nil
}

func (r *MemoryRepository) Get(key domain.GoroutineId) (string, error) {
	if value, exists := r.values[key]; exists {
		return value, nil
	}
	return "", fmt.Errorf("key not found: %d", key)
}

func (r *MemoryRepository) Delete(key domain.GoroutineId) error {
	delete(r.values, key)
	return nil
}
