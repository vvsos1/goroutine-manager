package domain

import "time"

type Goroutine struct {
	ID           int64
	SaveDuration time.Duration
	repository   KeyValueRepository[string, string]
	channel      chan struct{}
	started      bool
}

func (g *Goroutine) Start() {
	if g.started {
		return
	}
}

var nextId int64 = 1

func nextGoroutineId() int64 {
	nextId++
	return nextId
}

func NewGoroutine(duration time.Duration, repository KeyValueRepository[string, string]) Goroutine {
	return Goroutine{
		ID:           nextGoroutineId(),
		SaveDuration: duration,
		repository:   repository,
		channel:      make(chan struct{}),
		started:      false,
	}
}
