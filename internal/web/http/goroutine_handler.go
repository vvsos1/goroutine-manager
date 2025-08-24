package router

import (
	"encoding/json"
	"goroutine-manager/internal/domain"
	"goroutine-manager/internal/usecase"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type GoroutineHandler struct {
	usecase usecase.GoroutineUsecase
}

func NewGoroutineHandler(usecase usecase.GoroutineUsecase) *GoroutineHandler {
	return &GoroutineHandler{usecase: usecase}
}

func (h *GoroutineHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.createGoroutine)
	r.Get("/", h.countGoroutines)
	r.Get("/{id}", h.getGoroutine)
	r.Patch("/{id}", h.updateGoroutine)
	r.Delete("/{id}", h.deleteGoroutine)
	return r
}

func (h *GoroutineHandler) createGoroutine(w http.ResponseWriter, r *http.Request) {
	var body struct {
		SaveDuration int `json:"save_duration"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Cannot Decode Body", http.StatusBadRequest)
		return
	}
	id, err := h.usecase.Create(body.SaveDuration)
	if err != nil {
		http.Error(w, "Cannot Create Goroutine", http.StatusInternalServerError)
		return
	}

	writeJson(w, http.StatusCreated, map[string]string{"id": strconv.Itoa(int(id))})

}

func (h *GoroutineHandler) getGoroutine(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	value, err := h.usecase.Get(domain.GoroutineId(intId))
	if err != nil {
		http.Error(w, "Goroutine Not Found", http.StatusNotFound)
		return
	}

	writeJson(w, http.StatusOK, map[string]string{"id": strconv.Itoa(intId), "value": value})
}

func (h *GoroutineHandler) deleteGoroutine(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	err = h.usecase.Delete(domain.GoroutineId(intId))
	if err != nil {
		http.Error(w, "Goroutine Not Found", http.StatusNotFound)
		return
	}

	writeJson(w, http.StatusNoContent, nil)
}

func (h *GoroutineHandler) countGoroutines(w http.ResponseWriter, r *http.Request) {
	count := h.usecase.Count()
	writeJson(w, http.StatusOK, map[string]int{"count": count})
}

func (h *GoroutineHandler) updateGoroutine(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var body struct {
		SaveDuration int `json:"save_duration"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Cannot Decode Body", http.StatusBadRequest)
		return
	}

	err = h.usecase.Update(domain.GoroutineId(intId), body.SaveDuration)
	if err != nil {
		http.Error(w, "Goroutine Not Found", http.StatusNotFound)
		return
	}

	writeJson(w, http.StatusNoContent, nil)
}

func writeJson(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}
