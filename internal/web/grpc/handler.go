package grpc

import (
	"context"
	"fmt"
	"goroutine-manager/internal/domain"
	"goroutine-manager/internal/usecase"
	pb "goroutine-manager/internal/web/grpc/pb/worker"
)

type WorkerHandler struct {
	pb.UnimplementedWorkerServiceServer
	usecase usecase.WorkerUsecase
}

func NewWorkerHandler(usecase usecase.WorkerUsecase) *WorkerHandler {
	return &WorkerHandler{
		usecase: usecase,
	}
}

func (w *WorkerHandler) CreateWorker(_ context.Context, r *pb.CreateWorkerRequest) (*pb.CreateWorkerResponse, error) {
	workerId, err := w.usecase.Create(int(r.SaveDuration), r.WorkerMsg)
	if err != nil {
		return nil, err
	}
	return &pb.CreateWorkerResponse{WorkerId: int64(workerId), Msg: "worker successfully created"}, nil
}

func (w *WorkerHandler) CountWorkers(_ context.Context, _ *pb.CountWorkersRequest) (*pb.CountWorkersResponse, error) {
	count := w.usecase.Count()
	return &pb.CountWorkersResponse{WorkerCount: int64(count)}, nil
}

func (w *WorkerHandler) GetWorker(_ context.Context, r *pb.GetWorkerRequest) (*pb.GetWorkerResponse, error) {
	worker, err := w.usecase.Get(domain.WorkerId(r.WorkerId))
	if err != nil {
		return nil, fmt.Errorf("failed to get worker: %w", err)
	}
	return &pb.GetWorkerResponse{
		WorkerId:  int64(worker.Id),
		WorkerMsg: worker.WorkerMsg,
		Status:    string(worker.Status),
	}, nil
}

func (w *WorkerHandler) GetWorkerData(_ context.Context, r *pb.GetWorkerDataRequest) (*pb.GetWorkerDataResponse, error) {
	data, err := w.usecase.GetData(domain.WorkerId(r.WorkerId))
	if err != nil {
		return nil, fmt.Errorf("failed to get worker data: %w", err)
	}
	return &pb.GetWorkerDataResponse{
		LastModified: data.LastModified.String(),
		CachedMsg:    data.WorkerMsg,
	}, nil
}

func (w *WorkerHandler) UpdateWorker(_ context.Context, r *pb.UpdateWorkerRequest) (*pb.UpdateWorkerResponse, error) {
	err := w.usecase.Update(domain.WorkerId(r.WorkerId), int(r.SaveDuration), r.WorkerMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to update worker: %w", err)
	}
	return &pb.UpdateWorkerResponse{Msg: "worker successfully updated"}, nil
}

func (w *WorkerHandler) DeleteWorker(_ context.Context, r *pb.DeleteWorkerRequest) (*pb.DeleteWorkerResponse, error) {
	err := w.usecase.Delete(domain.WorkerId(r.WorkerId))
	if err != nil {
		return nil, fmt.Errorf("failed to delete worker: %w", err)
	}
	return &pb.DeleteWorkerResponse{Msg: "worker successfully deleted"}, nil
}
