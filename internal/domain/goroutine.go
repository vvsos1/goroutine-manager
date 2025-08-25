package domain

import (
	"log"
	"time"
)

type GoroutineId int

type Goroutine struct {
	Id               GoroutineId
	SaveDuration     int
	dataRepository   KeyValueRepository
	operationChannel chan interface{}
	resultChannel    chan string
	started          bool
}

type GoroutineRepository interface {
	Save(goroutine *Goroutine) error
	Get(id GoroutineId) (*Goroutine, error)
	Delete(id GoroutineId) error
	Count() int
}

func (g *Goroutine) StartInGoroutine() {
	if g.started {
		return
	}
	go g.process()
}

func (g *Goroutine) process() {
	log.Println("Goroutine started with ID:", g.Id)
	g.started = true
	defer func() {
		g.started = false
		log.Println("Goroutine ended with ID:", g.Id)
	}()
	tick := time.Tick(time.Second * time.Duration(g.SaveDuration))
	for {
		select {
		case operation := <-g.operationChannel:
			switch op := operation.(type) {
			case *ReadOperation:
				value := g.readFromRepository()
				g.resultChannel <- value
			case *DeleteOperation:
				return
			case *UpdateOperation:
				g.SaveDuration = op.SaveDuration
				tick = time.Tick(time.Second * time.Duration(g.SaveDuration))
			}
		case <-tick:
			g.saveToRepository()
		}
	}
}

func (g *Goroutine) Read() string {
	operation := &ReadOperation{}
	g.operationChannel <- operation
	value := <-g.resultChannel
	return value
}

func (g *Goroutine) Update(saveDuration int) {
	operation := &UpdateOperation{SaveDuration: saveDuration}
	g.operationChannel <- operation
}

func (g *Goroutine) Delete() {
	operation := &DeleteOperation{}
	g.operationChannel <- operation
	close(g.operationChannel) // 채널을 닫아 고루틴이 종료되도록 함
}

func (g *Goroutine) saveToRepository() {
	err := g.dataRepository.Put(g.Id, time.Now().String())
	if err != nil {
		log.Println("Failed to save to data repository:", err)
		return
	}
}

func (g *Goroutine) readFromRepository() string {
	value, err := g.dataRepository.Get(g.Id)
	if err != nil {
		log.Println("Failed to read from data repository:", err)
		return ""
	}
	return value
}

var nextId GoroutineId = 0

func nextGoroutineId() GoroutineId {
	nextId++
	return nextId
}

func NewGoroutine(saveDuration int, repository KeyValueRepository) *Goroutine {
	return &Goroutine{
		Id:             nextGoroutineId(),
		SaveDuration:   saveDuration,
		dataRepository: repository,
		// 고루틴에게 명령을 전달하는 채널
		operationChannel: make(chan interface{}),
		// 고루틴으로부터 결과를 전달받는 채널
		resultChannel: make(chan string),
		started:       false,
	}
}

type UpdateOperation struct {
	SaveDuration int
}

type ReadOperation struct {
}

type DeleteOperation struct {
}
