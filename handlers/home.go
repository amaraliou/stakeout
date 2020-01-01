package handlers

import (
	"net/http"

	"github.com/amaraliou/apetitoso/responses"
)

// Home -> handles GET /api/v1/
func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To This Awesome API")
}
