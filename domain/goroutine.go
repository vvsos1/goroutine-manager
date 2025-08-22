package domain

import (
	"time"
)

type GoroutineId int

type Goroutine struct {
	Id           GoroutineId
	SaveDuration time.Duration
	repository   KeyValueRepository[GoroutineId, string]
	channel      chan interface{}
	started      bool
}

func (g *Goroutine) StartInGoroutine() {
	if g.started {
		return
	}
	go g.process()
}

func (g *Goroutine) process() {
	g.started = true
	defer func() { g.started = false }()
	tick := time.Tick(g.SaveDuration)
	for {
		select {
		case operation := <-g.channel:
			switch op := operation.(type) {
			case *ReadOperation:
				value := g.readFromRepository()
				g.channel <- value
			case *DeleteOperation:
				return
			case *UpdateOperation:
				g.SaveDuration = op.SaveDuration
				tick = time.Tick(g.SaveDuration)
			}
		case <-tick:
			g.saveToRepository()
		}
	}
}

func (g *Goroutine) Read() string {
	operation := &ReadOperation{}
	g.channel <- operation
	value := <-g.channel
	if strValue, ok := value.(string); ok {
		return strValue
	}
	// TODO: 에러 처리
	return ""
}

func (g *Goroutine) Update(duration time.Duration) {
	operation := &UpdateOperation{SaveDuration: duration}
	g.channel <- operation
}

func (g *Goroutine) Delete() {
	operation := &DeleteOperation{}
	g.channel <- operation
	close(g.channel) // 채널을 닫아 고루틴이 종료되도록 함
}

func (g *Goroutine) saveToRepository() {
	err := g.repository.Put(g.Id, time.Now().String())
	if err != nil {
		// TODO: 에러 처리
		return
	}
}

func (g *Goroutine) readFromRepository() string {
	value, err := g.repository.Get(g.Id)
	if err != nil {
		// TODO: 에러 처리
		return ""
	}
	return value
}

var nextId GoroutineId = 1

func nextGoroutineId() GoroutineId {
	nextId++
	return nextId
}

func NewGoroutine(duration time.Duration, repository KeyValueRepository[GoroutineId, string]) Goroutine {
	return Goroutine{
		Id:           nextGoroutineId(),
		SaveDuration: duration,
		repository:   repository,
		channel:      make(chan interface{}),
		started:      false,
	}
}

type UpdateOperation struct {
	SaveDuration time.Duration
}

type ReadOperation struct {
}

type DeleteOperation struct {
}
