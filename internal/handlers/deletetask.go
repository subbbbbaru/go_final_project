package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/subbbbbaru/first-sample/pkg/log"
)

func (h *Handler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if len(id) == 0 {
		log.Error().Println("wrong id")
		w.WriteHeader(http.StatusBadRequest)
		if errJson := json.NewEncoder(w).Encode(map[string]string{"error": "wrong id"}); errJson != nil {
			log.Error().Println(errJson)
		}
		return
	}
	taskId, err := strconv.Atoi(id)
	if err != nil {
		log.Error().Println(err)
		w.WriteHeader(http.StatusBadRequest)
		if errJson := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); errJson != nil {
			log.Error().Println(errJson)
		}
		return
	}

	if errDelete := h.services.Delete(taskId); errDelete != nil {
		log.Error().Println(errDelete)
		w.WriteHeader(http.StatusInternalServerError)
		if errJson := json.NewEncoder(w).Encode(map[string]string{"error": errDelete.Error()}); errJson != nil {
			log.Error().Println(errJson)
		}
		return
	}

	w.Header().Set("Content-type", "application-json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if errJson := json.NewEncoder(w).Encode(map[string]string{}); errJson != nil {
		log.Error().Println(errJson)
	}
}
