package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/subbbbbaru/first-sample/pkg/log"

	"github.com/subbbbbaru/go_final_project/internal/models"
)

func (h *Handler) CrateTaskHandler(w http.ResponseWriter, r *http.Request) {

	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Error().Println(err)
		w.WriteHeader(http.StatusBadRequest)
		if errJson := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); errJson != nil {
			log.Error().Println(errJson)
		}
		return
	}

	id, err := h.services.TodoTask.Create(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if errJson := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); errJson != nil {
			log.Error().Println(errJson)
		}
		return
	}

	response := struct {
		ID int64 `json:"id"`
	}{ID: id}

	w.Header().Set("Content-type", "application-json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if errJson := json.NewEncoder(w).Encode(response); errJson != nil {
		log.Error().Println(errJson)
	}
}
