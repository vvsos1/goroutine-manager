package grpc

import (
	"context"
	"fmt"
	pb "worker-manager/api/worker"
	"worker-manager/internal/domain"
	"worker-manager/internal/usecase"
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

func (w *WorkerHandler) CreateWorker(ctx context.Context, r *pb.CreateWorkerRequest) (*pb.CreateWorkerResponse, error) {
	workerId, err := w.usecase.Create(ctx, int(r.SaveDuration), r.WorkerMsg)
	if err != nil {
		return nil, err
	}
	return &pb.CreateWorkerResponse{WorkerId: int64(workerId), Msg: "worker successfully created"}, nil
}

func (w *WorkerHandler) CountWorkers(ctx context.Context, _ *pb.CountWorkersRequest) (*pb.CountWorkersResponse, error) {
	count := w.usecase.Count(ctx)
	return &pb.CountWorkersResponse{WorkerCount: int64(count)}, nil
}

func (w *WorkerHandler) GetWorker(ctx context.Context, r *pb.GetWorkerRequest) (*pb.GetWorkerResponse, error) {
	worker, err := w.usecase.Get(ctx, domain.WorkerId(r.WorkerId))
	if err != nil {
		return nil, fmt.Errorf("failed to get worker: %w", err)
	}
	return &pb.GetWorkerResponse{
		WorkerId:  int64(worker.Id),
		WorkerMsg: worker.WorkerMsg,
		Status:    string(worker.Status),
	}, nil
}

func (w *WorkerHandler) GetWorkerData(ctx context.Context, r *pb.GetWorkerDataRequest) (*pb.GetWorkerDataResponse, error) {
	data, err := w.usecase.GetData(ctx, domain.WorkerId(r.WorkerId))
	if err != nil {
		return nil, fmt.Errorf("failed to get worker data: %w", err)
	}
	return &pb.GetWorkerDataResponse{
		LastModified: data.LastModified.String(),
		CachedMsg:    data.WorkerMsg,
	}, nil
}

func (w *WorkerHandler) UpdateWorker(ctx context.Context, r *pb.UpdateWorkerRequest) (*pb.UpdateWorkerResponse, error) {
	err := w.usecase.Update(ctx, domain.WorkerId(r.WorkerId), int(r.SaveDuration), r.WorkerMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to update worker: %w", err)
	}
	return &pb.UpdateWorkerResponse{Msg: "worker successfully updated"}, nil
}

func (w *WorkerHandler) DeleteWorker(ctx context.Context, r *pb.DeleteWorkerRequest) (*pb.DeleteWorkerResponse, error) {
	err := w.usecase.Delete(ctx, domain.WorkerId(r.WorkerId))
	if err != nil {
		return nil, fmt.Errorf("failed to delete worker: %w", err)
	}
	return &pb.DeleteWorkerResponse{Msg: "worker successfully deleted"}, nil
}
