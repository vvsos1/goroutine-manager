package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	pb "worker-manager/api/worker"
	"worker-manager/internal/domain"
	"worker-manager/internal/infra/repository/data"
	"worker-manager/internal/infra/repository/worker"
	"worker-manager/internal/usecase"
	"worker-manager/internal/web/grpc"
	router "worker-manager/internal/web/http"
	grpcMiddleware "worker-manager/middleware/grpc"
	httpMiddleware "worker-manager/middleware/http"
	"worker-manager/util/logger"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	gogrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	//dataRepo := data.NewMemoryRepository()
	useMemoryRepo := os.Getenv("USE_MEMORY_REPO")
	var dataRepo domain.DataRepository

	if useMemoryRepo == "true" {
		logger.Infoln(context.Background(), "Using in-memory data repository")
		dataRepo = data.NewMemoryRepository()
	} else {
		logger.Infoln(context.Background(), "Using Valkey data repository")
		addr := os.Getenv("VALKEY_ADDR")
		if addr == "" {
			addr = "localhost:6379"
		}
		dataRepo = data.NewValkeyRepository(addr)
	}
	workerRepository := worker.NewMemoryWorkerRepository()
	workerUsecase := usecase.NewWorkerInteractor(workerRepository, dataRepo)

	errs := make(chan error)

	go func() {
		r := chi.NewRouter()
		r.Use(httpMiddleware.RequestId)
		r.Use(chiMiddleware.Logger)
		r.Use(chiMiddleware.Recoverer)

		router.MountRoutes(r, workerUsecase)
		httpPort := os.Getenv("HTTP_PORT")
		if httpPort == "" {
			httpPort = "3000"
		}
		logger.Infoln(context.Background(), "HTTP server listening on :", httpPort)

		errs <- http.ListenAndServe(":"+httpPort, r)
	}()

	go func() {
		grpcPort := os.Getenv("GRPC_PORT")
		if grpcPort == "" {
			grpcPort = "3001"
		}
		listen, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			errs <- err
		}

		handler := grpc.NewWorkerHandler(workerUsecase)
		server := gogrpc.NewServer(
			gogrpc.UnaryInterceptor(
				grpcMiddleware.RequestId))
		pb.RegisterWorkerServiceServer(server, handler)
		// goland grpc call test용 리플렉션 추가
		reflection.Register(server)
		logger.Infoln(context.Background(), "GRPC server listening on :", grpcPort)
		errs <- server.Serve(listen)
	}()

	go func() {
		// Graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-sigChan)
	}()

	logger.Errorf(context.Background(), "terminated %s", <-errs)
}
