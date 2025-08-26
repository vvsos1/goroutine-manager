package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"worker-manager/internal/domain"
	"worker-manager/internal/infra/repository/data"
	"worker-manager/internal/infra/repository/worker"
	"worker-manager/internal/usecase"
	"worker-manager/internal/web/grpc"
	pb "worker-manager/internal/web/grpc/pb/worker"
	router "worker-manager/internal/web/http"
	middleware2 "worker-manager/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	gogrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	//dataRepo := data.NewMemoryRepository()
	useMemoryRepo := os.Getenv("USE_MEMORY_REPO")
	var dataRepo domain.DataRepository

	if useMemoryRepo == "true" {
		log.Println("Using in-memory data repository")
		dataRepo = data.NewMemoryRepository()
	} else {
		log.Println("Using Valkey data repository")
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
		r.Use(middleware2.SetRequestIdMiddleware)
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)

		router.MountRoutes(r, workerUsecase)
		httpPort := os.Getenv("HTTP_PORT")
		if httpPort == "" {
			httpPort = "3000"
		}
		log.Printf("HTTP server listening on :%s", httpPort)

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
		server := gogrpc.NewServer()
		pb.RegisterWorkerServiceServer(server, handler)
		reflection.Register(server)
		log.Printf("GRPC server listening on :%s", grpcPort)
		errs <- server.Serve(listen)
	}()

	go func() {
		// Graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-sigChan)
	}()

	log.Fatalf("terminated %s", <-errs)
}
