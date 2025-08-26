package domain

import (
	"context"
	"time"
)

type Data struct {
	WorkerId     WorkerId  `json:"worker_id"`
	WorkerMsg    string    `json:"worker_msg"`
	LastModified time.Time `json:"last_modified"`
}

type DataRepository interface {
	Put(ctx context.Context, key WorkerId, value *Data) error

	Get(ctx context.Context, key WorkerId) (*Data, error)

	Delete(ctx context.Context, key WorkerId) error
}
