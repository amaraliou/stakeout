package middlewares

import (
	"errors"
	"net/http"

	"github.com/amaraliou/apetitoso/auth"
	"github.com/amaraliou/apetitoso/responses"
)

// SetMiddlewareAuthentication ...
func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		next(w, r)
	}
}
