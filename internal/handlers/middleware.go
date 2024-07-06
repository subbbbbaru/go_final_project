package handlers

import (
	"net/http"

	"github.com/subbbbbaru/first-sample/pkg/log"
)

func (h *Handler) UserIdentity(next http.HandlerFunc, password string) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if len(password) > 0 {
				cookie, err := r.Cookie("token")
				if err != nil {
					http.Error(w, "Authentificate required", http.StatusUnauthorized)

					log.Error().Println("Authentificate required", err.Error())
					return
				}
				valid, err := h.services.Auth.ValideToken(cookie.Value)
				if err != nil {
					http.Error(w, "Internal server error", http.StatusUnauthorized)
					log.Error().Println("Internal server error", err.Error())
					return
				}
				if !valid {
					http.Error(w, "Authentificate required", http.StatusUnauthorized)
					log.Error().Println("Authentificate required")
					return
				}
			}
			next(w, r)
		},
	)
}
