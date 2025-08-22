package domain

import "time"

type GoroutineManager struct {
	goroutines map[int64]Goroutine
	repository KeyValueRepository[string, string]
}

func NewGoroutineManager() *GoroutineManager {
	return &GoroutineManager{
		goroutines: make(map[int64]Goroutine),
	}
}

func (gm *GoroutineManager) Create(req CreateGoroutine) (Goroutine, error) {
	goroutine := NewGoroutine(req.SaveDuration)
	gm.goroutines[goroutine.ID] = goroutine
	return goroutine, nil
}

// 고루틴 생성 DTO
type CreateGoroutine struct {
	SaveDuration time.Duration `json:"duration"`
}
