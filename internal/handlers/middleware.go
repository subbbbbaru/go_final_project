package handlers

import (
	"log"
	"net/http"
)

func (h *Handler) UserIdentity(next http.HandlerFunc, password string) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if len(password) > 0 {
				cookie, err := r.Cookie("token")
				if err != nil {
					http.Error(w, "Authentificate required", http.StatusUnauthorized)
					log.Println("Authentificate required", err.Error())
					return
				}
				valid, err := h.services.Auth.ValideToken(cookie.Value)
				if err != nil {
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					log.Println("Internal server error", err.Error())
					return
				}
				if !valid {
					http.Error(w, "Authentificate required", http.StatusUnauthorized)
					log.Println("Authentificate required")
					return
				}
			}
			next(w, r)
		},
	)
}

// 	header := c.GetHeader(authHeader)
// 	if header == "" {
// 		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
// 		return
// 	}
// 	headerParts := strings.Split(header, " ")
// 	if len(headerParts) != 2 {
// 		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
// 		return
// 	}

// 	userId, err := h.services.Auth.ParseToken(headerParts[1])
// 	if err != nil {
// 		newErrorResponse(c, http.StatusUnauthorized, err.Error())
// 	}
// 	c.Set(userCtx, userId)
// }
