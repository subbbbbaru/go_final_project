package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/subbbbbaru/first-sample/pkg/log"
)

func (h *Handler) DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusInternalServerError)
		if errJson := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); errJson != nil {
			log.Error().Println(errJson)
		}
		return
	}

	if errDone := h.services.Done(taskId); errDone != nil {
		log.Error().Println(errDone)
		w.WriteHeader(http.StatusInternalServerError)
		if errJson := json.NewEncoder(w).Encode(map[string]string{"error": errDone.Error()}); errJson != nil {
			log.Error().Println(errJson)
		}
		return
	}

	w.Header().Set("Content-type", "application-json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]string{})
	if err != nil {
		log.Error().Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		if errJson := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); errJson != nil {
			log.Error().Println(errJson)
		}
		return
	}
}
