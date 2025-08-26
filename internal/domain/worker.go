package domain

import (
	"context"
	"time"
	"worker-manager/util/logger"
)

type WorkerId int

type WorkerStatus string

const (
	UP   WorkerStatus = "UP"
	DOWN WorkerStatus = "DOWN"
)

type Worker struct {
	Id           WorkerId
	SaveDuration int
	WorkerMsg    string
	// 데이터를 저장할 레포지토리
	dataRepository DataRepository
	// 고루틴에게 명령을 전달하는 채널
	operationChannel chan interface{}
	// 고루틴으로부터 결과를 전달받는 채널
	resultChannel chan *Data
	Status        WorkerStatus
}

type WorkerRepository interface {
	Save(ctx context.Context, worker *Worker) error
	Get(ctx context.Context, id WorkerId) (*Worker, error)
	//Delete(id WorkerId) error
	Count(ctx context.Context) int
}

func NewWorker(saveDuration int, workerMsg string, repository DataRepository) *Worker {
	return &Worker{
		Id:             nextWorkerId(),
		SaveDuration:   saveDuration,
		WorkerMsg:      workerMsg,
		dataRepository: repository,
		// 고루틴에게 명령을 전달하는 채널
		operationChannel: make(chan interface{}),
		// 고루틴으로부터 결과를 전달받는 채널
		resultChannel: make(chan *Data),
		Status:        DOWN,
	}
}

func (w *Worker) StartInGoroutine() {
	if w.Status == UP {
		return
	}
	go w.process()
}

func (w *Worker) process() {
	logger.Debugln(context.Background(), "Worker started with ID:", w.Id)
	w.Status = UP
	defer func() {
		w.Status = DOWN
		logger.Debugln(context.Background(), "Worker ended with ID:", w.Id)
	}()
	tick := time.Tick(time.Second * time.Duration(w.SaveDuration))
	for {
		select {
		case operation := <-w.operationChannel:
			switch op := operation.(type) {
			case *readDataOperation:
				value := w.readFromRepository()
				w.resultChannel <- value
			case *deleteOperation:
				w.deleteFromRepository()
				return
			case *updateOperation:
				w.SaveDuration = op.SaveDuration
				w.WorkerMsg = op.WorkerMsg
				tick = time.Tick(time.Second * time.Duration(w.SaveDuration))
			}
		case <-tick:
			w.saveToRepository()
		}
	}
}

func (w *Worker) ReadData() *Data {
	operation := &readDataOperation{}
	w.operationChannel <- operation
	value := <-w.resultChannel
	return value
}

func (w *Worker) Read() *Data {
	return w.readFromRepository()
}

func (w *Worker) Update(saveDuration int, workerMsg string) {
	operation := &updateOperation{SaveDuration: saveDuration, WorkerMsg: workerMsg}
	w.operationChannel <- operation
}

func (w *Worker) Delete() {
	operation := &deleteOperation{}
	w.operationChannel <- operation
	close(w.operationChannel) // 채널을 닫아 고루틴이 종료되도록 함]
	close(w.resultChannel)
}

func (w *Worker) saveToRepository() {
	data := &Data{
		WorkerId:     w.Id,
		WorkerMsg:    w.WorkerMsg,
		LastModified: time.Now(),
	}
	err := w.dataRepository.Put(context.Background(), w.Id, data)
	if err != nil {
		logger.Infoln(context.Background(), "Failed to save to data repository:", err)
		return
	}
}

func (w *Worker) readFromRepository() *Data {
	value, err := w.dataRepository.Get(context.Background(), w.Id)
	if err != nil {
		logger.Infoln(context.Background(), "Failed to read from data repository:", err)
		return nil
	}
	return value
}

func (w *Worker) deleteFromRepository() {
	err := w.dataRepository.Delete(context.Background(), w.Id)
	if err != nil {
		logger.Infoln(context.Background(), "Failed to delete from data repository:", err)
	}
}

var nextId WorkerId = 0

func nextWorkerId() WorkerId {
	nextId++
	return nextId
}

type updateOperation struct {
	SaveDuration int
	WorkerMsg    string
}

type readDataOperation struct {
}

type deleteOperation struct {
}
