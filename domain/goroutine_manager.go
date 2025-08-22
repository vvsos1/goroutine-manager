package domain

import (
	"fmt"
	"time"
)

type GoroutineManager struct {
	goroutines map[GoroutineId]Goroutine
	repository KeyValueRepository[GoroutineId, string]
}

func NewGoroutineManager() *GoroutineManager {
	return &GoroutineManager{
		goroutines: make(map[GoroutineId]Goroutine),
	}
}

func (gm *GoroutineManager) Create(req CreateGoroutine) (Goroutine, error) {
	goroutine := NewGoroutine(req.SaveDuration, gm.repository)
	gm.goroutines[goroutine.Id] = goroutine

	goroutine.StartInGoroutine()

	return goroutine, nil
}

func (gm *GoroutineManager) Get(id GoroutineId) (string, error) {
	goroutine, exists := gm.goroutines[id]
	if !exists {
		return "", fmt.Errorf("goroutine with id %d not found", id)
	}
	value := goroutine.Read()
	return value, nil
}

func (gm *GoroutineManager) Update(id GoroutineId, duration time.Duration) error {
	goroutine, exists := gm.goroutines[id]
	if !exists {
		return fmt.Errorf("goroutine with id %d not found", id)
	}
	goroutine.Update(duration)
	return nil
}

func (gm *GoroutineManager) Delete(id GoroutineId) {
	goroutine, exists := gm.goroutines[id]
	if !exists {
	}
	goroutine.Delete()
	delete(gm.goroutines, id)
}

func (gm *GoroutineManager) Count() int {
	return len(gm.goroutines)
}

// 고루틴 생성 DTO
type CreateGoroutine struct {
	SaveDuration time.Duration `json:"duration"`
}
