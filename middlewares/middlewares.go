package middlewares

import (
	"errors"
	"github.com/stepan-romankov/go-by-bike/auth"
	"github.com/stepan-romankov/go-by-bike/responses"
	"log"
	"net/http"
)

type AuthMiddleware struct {
	AuthStore auth.Store
}

func NewAuthMiddleware(authStore auth.Store) *AuthMiddleware {
	return &AuthMiddleware{AuthStore: authStore}
}

func (am *AuthMiddleware) Wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := am.AuthStore.GetUserId(r)
		if err != nil {
			log.Printf("Unauthorized access")
			responses.Error(w, http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}

		next(w, r)
	}
}
