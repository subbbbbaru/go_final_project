package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/subbbbbaru/first-sample/pkg/log"
	"github.com/subbbbbaru/go_final_project/internal/models"
	"github.com/subbbbbaru/go_final_project/internal/service"
	"github.com/subbbbbaru/go_final_project/utils"
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
		log.Info().Println("method not yet realized")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not yet realized"})
		return
	}
}

func (h *Handler) CrateTaskHandler(w http.ResponseWriter, r *http.Request) {

	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Error().Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if task.Title == "" {
		log.Error().Println("task title not found")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "task title not found"})
		return
	}

	timeNow := time.Now() //.Truncate(24 * time.Hour)

	if task.Date == "" {
		task.Date = timeNow.Format("20060102")
	}

	date, err := time.Parse("20060102", task.Date)
	if err != nil {
		log.Error().Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if date.Before(timeNow) {
		if task.Repeat == "" {
			task.Date = timeNow.Format("20060102")
		} else {
			nextDate, err := utils.NextDate(timeNow, task.Date, task.Repeat)
			if err != nil {
				log.Error().Println(err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			task.Date = nextDate
		}
	}

	id, err := h.services.TodoTask.Create(task)
	if err != nil {
		log.Error().Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	response := struct {
		ID int64 `json:"id"`
	}{ID: id}

	w.Header().Set("Content-type", "application-json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
