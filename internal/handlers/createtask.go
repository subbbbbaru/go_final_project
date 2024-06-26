package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/subbbbbaru/go_final_project/internal/models"
	"github.com/subbbbbaru/go_final_project/internal/repository"
	"github.com/subbbbbaru/go_final_project/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CrateTaskHandler(w, r)
	case http.MethodGet:
		h.GetTaskByIdHandler(w, r)
	case http.MethodPut:
		h.PutTaskHandler(w, r)
	case http.MethodDelete:
		h.DeleteTaskHandler(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not yet realized"})
		return
	}
}

func (h *Handler) CrateTaskHandler(w http.ResponseWriter, r *http.Request) {

	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "task title not found"})
		return
	}

	timeNow := time.Now().Truncate(24 * time.Hour).UTC()

	if task.Date == "" {
		task.Date = timeNow.Format("20060102")
	}

	date, err := time.Parse("20060102", task.Date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	repos := repository.NewRepository(nil)
	if date.Before(timeNow) {
		if task.Repeat == "" {
			task.Date = timeNow.Format("20060102")
		} else {
			task.Date, err = repos.NextDate(timeNow, task.Date, task.Repeat)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
		}
	}

	id, err := h.services.TodoTask.Create(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	response := struct {
		ID int64 `json:"id"`
	}{ID: id}

	// jsonee, err := json.Marshal(response)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	// 	return
	// }
	w.Header().Set("Content-type", "application-json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
