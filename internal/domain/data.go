package domain

import (
	"time"
)

type Data struct {
	WorkerId     WorkerId  `json:"worker_id"`
	WorkerMsg    string    `json:"worker_msg"`
	LastModified time.Time `json:"last_modified"`
}

type DataRepository interface {
	Put(key WorkerId, value *Data) error

	Get(key WorkerId) (*Data, error)

	Delete(key WorkerId) error
}
