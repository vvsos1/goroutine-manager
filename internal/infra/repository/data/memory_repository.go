package data

import "goroutine-manager/internal/domain"

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
	return "", nil // Return zero value of V if key does not exist
}
