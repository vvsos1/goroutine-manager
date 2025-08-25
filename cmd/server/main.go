package main

import (
	"goroutine-manager/internal/domain"
	"goroutine-manager/internal/infra/repository/data"
	"goroutine-manager/internal/infra/repository/goroutine"
	"goroutine-manager/internal/usecase"
	router "goroutine-manager/internal/web/http"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	//dataRepo := data.NewMemoryRepository()
	useMemoryRepo := os.Getenv("USE_MEMORY_REPO")
	var dataRepo domain.KeyValueRepository

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
	goroutineRepo := goroutine.NewMemoryGoroutineRepository()
	goroutineUC := usecase.NewGoroutineInteractor(goroutineRepo, dataRepo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	router.MountRoutes(r, goroutineUC)

	// Read server port from env
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	listenAddr := port
	if !strings.Contains(port, ":") {
		listenAddr = ":" + port
	}
	log.Printf("HTTP server listening on %s", listenAddr)

	if err := http.ListenAndServe(listenAddr, r); err != nil {
		panic(err)
	}

}
