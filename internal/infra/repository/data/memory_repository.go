package data

import (
	"context"
	"fmt"
	"worker-manager/internal/domain"
	"worker-manager/util/logger"
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

func (r *MemoryRepository) Put(ctx context.Context, key domain.WorkerId, value *domain.Data) error {
	r.values[key] = value
	logger.Infof(ctx, "memory repository put key: %d, value: %+v", key, value)
	return nil
}

func (r *MemoryRepository) Get(ctx context.Context, key domain.WorkerId) (*domain.Data, error) {
	if value, exists := r.values[key]; exists {
		logger.Infof(ctx, "memory repository get key: %d, value: %+v", key, value)
		return value, nil
	}
	logger.Infof(ctx, "memory repository get key: %d, value: nil", key)
	return nil, fmt.Errorf("key not found: %d", key)
}

func (r *MemoryRepository) Delete(ctx context.Context, key domain.WorkerId) error {
	delete(r.values, key)
	logger.Infof(ctx, "memory repository delete key: %d", key)
	return nil
}
