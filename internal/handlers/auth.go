package handlers

import (
	"encoding/json"

	"net/http"

	"github.com/subbbbbaru/first-sample/pkg/log"
)

type signInInput struct {
	Password string `json:"password" binding:"required"`
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var input signInInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Error().Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	token, err := h.services.Auth.GenerateToken(input.Password)
	if err != nil {
		log.Error().Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	response := struct {
		Token string `json:"token"`
	}{Token: token}

	w.Header().Set("Content-type", "application-json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
