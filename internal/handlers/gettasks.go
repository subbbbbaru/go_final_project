package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/subbbbbaru/first-sample/pkg/log"

	"github.com/subbbbbaru/go_final_project/internal/models"
)

func (h *Handler) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		if errJson := json.NewEncoder(w).Encode(map[string]string{"error": "method not yet realized"}); errJson != nil {
			log.Error().Println(errJson)
		}
		return
	}
	search := r.URL.Query().Get("search")

	tasks, err := h.services.TodoTask.GetTasks(search)
	if err != nil {
		log.Error().Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		if errJson := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); errJson != nil {
			log.Error().Println(errJson)
		}
		return
	}
	response := struct {
		Tasks []models.Task `json:"tasks"`
	}{Tasks: tasks}

	w.Header().Set("Content-type", "application-json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error().Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		if errJson := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); errJson != nil {
			log.Error().Println(errJson)
		}
		return
	}
}
