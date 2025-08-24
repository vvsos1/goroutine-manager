package router

import (
	"goroutine-manager/internal/usecase"

	"github.com/go-chi/chi/v5"
)

func MountRoutes(r chi.Router, usecase usecase.GoroutineUsecase) {
	handler := NewGoroutineHandler(usecase)
	r.Mount("/goroutines", handler.Routes())
}
