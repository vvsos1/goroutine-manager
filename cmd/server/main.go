package main

import (
	"goroutine-manager/internal/infra/repository/data"
	"goroutine-manager/internal/infra/repository/goroutine"
	"goroutine-manager/internal/usecase"
	router "goroutine-manager/internal/web/http"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	//dataRepo := data.NewMemoryRepository()
	dataRepo := data.NewValkeyRepository("localhost:6379")
	goroutineRepo := goroutine.NewMemoryGoroutineRepository()
	goroutineUC := usecase.NewGoroutineInteractor(goroutineRepo, dataRepo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	router.MountRoutes(r, goroutineUC)

	if err := http.ListenAndServe(":3000", r); err != nil {
		panic(err)
	}

}
