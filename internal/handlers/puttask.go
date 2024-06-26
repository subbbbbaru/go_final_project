package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/subbbbbaru/go_final_project/internal/models"
	"github.com/subbbbbaru/go_final_project/internal/repository"
)

func (h *Handler) PutTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if task.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("task ID = ", task.ID, " not found")
		json.NewEncoder(w).Encode(map[string]string{"error": "task ID not found"})
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

	ok, err := h.services.TodoTask.Update(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// jsonee, err := json.Marshal(response)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	// 	return
	// }
	w.Header().Set("Content-type", "application-json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ok)
}
