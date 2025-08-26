package router

import (
	"encoding/json"
	"net/http"
	"strconv"
	"worker-manager/internal/domain"
	"worker-manager/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type WorkerHandler struct {
	usecase usecase.WorkerUsecase
}

func NewWorkerHandler(usecase usecase.WorkerUsecase) *WorkerHandler {
	return &WorkerHandler{usecase: usecase}
}

func (h *WorkerHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.createWorker)
	r.Get("/", h.countWorkers)
	r.Get("/{id}/cache", h.getWorkerData)
	r.Get("/{id}", h.getWorker)
	r.Put("/{id}", h.updateWorker)
	r.Delete("/{id}", h.deleteWorker)
	return r
}

func (h *WorkerHandler) createWorker(w http.ResponseWriter, r *http.Request) {
	var body struct {
		SaveDuration int    `json:"save_duration"`
		WorkerMsg    string `json:"worker_msg"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Cannot Decode Body", http.StatusBadRequest)
		return
	}
	id, err := h.usecase.Create(r.Context(), body.SaveDuration, body.WorkerMsg)
	if err != nil {
		http.Error(w, "Cannot Create Worker", http.StatusInternalServerError)
		return
	}
	result := map[string]interface{}{
		"msg":       "worker successfully created",
		"worker_id": id,
	}
	writeJson(w, http.StatusCreated, result)

}

func (h *WorkerHandler) getWorkerData(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	value, err := h.usecase.GetData(r.Context(), domain.WorkerId(intId))
	if err != nil {
		http.Error(w, "Worker Not Found", http.StatusNotFound)
		return
	}

	result := map[string]interface{}{
		"last_modified": value.LastModified,
		"cached_msg":    value.WorkerMsg,
	}
	writeJson(w, http.StatusOK, result)
}

func (h *WorkerHandler) getWorker(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	value, err := h.usecase.Get(r.Context(), domain.WorkerId(intId))
	if err != nil {
		http.Error(w, "Worker Not Found", http.StatusNotFound)
		return
	}
	result := map[string]interface{}{
		"worker_id":  value.Id,
		"status":     value.Status,
		"worker_msg": value.WorkerMsg,
	}
	writeJson(w, http.StatusOK, result)
}

func (h *WorkerHandler) deleteWorker(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	err = h.usecase.Delete(r.Context(), domain.WorkerId(intId))
	if err != nil {
		http.Error(w, "Worker Not Found", http.StatusNotFound)
		return
	}
	result := map[string]interface{}{
		"msg": "worker successfully deleted",
	}
	writeJson(w, http.StatusOK, result)
}

func (h *WorkerHandler) countWorkers(w http.ResponseWriter, r *http.Request) {
	count := h.usecase.Count(r.Context())

	result := map[string]interface{}{
		"count": count,
	}
	writeJson(w, http.StatusOK, result)
}

func (h *WorkerHandler) updateWorker(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var body struct {
		SaveDuration int    `json:"save_duration"`
		WorkerMsg    string `json:"worker_msg"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Cannot Decode Body", http.StatusBadRequest)
		return
	}

	err = h.usecase.Update(r.Context(), domain.WorkerId(intId), body.SaveDuration, body.WorkerMsg)
	if err != nil {
		http.Error(w, "Worker Not Found", http.StatusNotFound)
		return
	}

	result := map[string]interface{}{
		"msg": "cache message changed",
	}
	writeJson(w, http.StatusOK, result)
}

func writeJson(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}
