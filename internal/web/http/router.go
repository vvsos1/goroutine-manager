package router

import (
	"goroutine-manager/internal/usecase"

	"github.com/go-chi/chi/v5"
)

func MountRoutes(r chi.Router, usecase usecase.WorkerUsecase) {
	handler := NewWorkerHandler(usecase)
	r.Mount("/workers", handler.Routes())
}
