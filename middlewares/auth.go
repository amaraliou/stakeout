package middlewares

import (
	"errors"
	"net/http"

	"github.com/amaraliou/stakeout/responses"
	"github.com/amaraliou/stakeout/auth"
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

// SetMiddlewareAdminAuthentication ...
func SetMiddlewareAdminAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}

		isAdmin, err := auth.IsAdminToken(r)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		if !isAdmin {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized: You are not an admin"))
			return
		}
		next(w, r)
	}
}
