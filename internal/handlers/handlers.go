package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/subbbbbaru/first-sample/pkg/log"

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
		log.Info().Println("method not yet realized")
		w.WriteHeader(http.StatusMethodNotAllowed)
		if errJson := json.NewEncoder(w).Encode(map[string]string{"error": "method not yet realized"}); errJson != nil {
			log.Error().Println(errJson)
		}
		return
	}
}
