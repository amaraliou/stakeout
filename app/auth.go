package app

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/amaraliou/apetitoso/models"
	"github.com/amaraliou/apetitoso/utils"
	"github.com/dgrijalva/jwt-go"
)

var JWTAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		notAuth := []string{"/api/user/new", "/api/user/login", "api/store/login", "api/store/new"}
		requestPath := request.URL.Path

		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(writer, request)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := request.Header.Get("Authorization")

		if tokenHeader == "" {
			response = utils.Message(false, "Missing auth token.")
			writer.WriteHeader(http.StatusForbidden)
			writer.Header().Add("Content-Type", "application/json")
			utils.Respond(writer, response)
			return
		}

		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			response = utils.Message(false, "Malformed/Invalid auth token.")
			writer.WriteHeader(http.StatusForbidden)
			writer.Header().Add("Content-Type", "application/json")
			utils.Respond(writer, response)
			return
		}

		tokenPart := splitted[1]
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN_PASSWORD")), nil
		})
		if err != nil {
			response = utils.Message(false, "Malformed auth token.")
			writer.WriteHeader(http.StatusForbidden)
			writer.Header().Add("Content-Type", "application/json")
			utils.Respond(writer, response)
			return
		}

		if !token.Valid {
			response = utils.Message(false, "Auth token is not valid.")
			writer.WriteHeader(http.StatusForbidden)
			writer.Header().Add("Content-Type", "application/json")
			utils.Respond(writer, response)
			return
		}

		ctx := context.WithValue(request.Context(), "user", tk.UserID)
		request = request.WithContext(ctx)
		next.ServeHTTP(writer, request)
	})
}
