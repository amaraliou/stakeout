package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/amaraliou/apetitoso/models"
	"github.com/amaraliou/apetitoso/responses"
	"github.com/gorilla/mux"
)

// CreateAdmin -> handles POST /api/v1/admin/
func (server *Server) CreateAdmin(writer http.ResponseWriter, request *http.Request) {

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
	}

	admin := models.Admin{}
	err = json.Unmarshal(body, &admin)
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	err = admin.Validate("create")
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	adminCreated, err := admin.CreateAdmin(server.DB)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	writer.Header().Set("Location", fmt.Sprintf("%s%s/%s", request.Host, request.RequestURI, adminCreated.ID.String()))
	responses.JSON(writer, http.StatusCreated, adminCreated)
}

// GetAdmins -> handles GET /api/v1/admin/
func (server *Server) GetAdmins(writer http.ResponseWriter, request *http.Request) {

	admin := models.Admin{}
	admins, err := admin.FindAllAdmins(server.DB)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(writer, http.StatusOK, admins)
}

// GetAdminByID -> handles GET /api/v1/admin/<id:uuid>
func (server *Server) GetAdminByID(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	admin := models.Admin{}
	adminRetrieved, err := admin.FindAdminByID(server.DB, vars["id"])
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(writer, http.StatusOK, adminRetrieved)
}

// UpdateAdmin -> handles PUT /api/v1/admin/<id:uuid>
func (server *Server) UpdateAdmin(writer http.ResponseWriter, request *http.Request) {

}

// DeleteAdmin -> handles DELETE /api/v1/admin/<id:uuid>
func (server *Server) DeleteAdmin(writer http.ResponseWriter, request *http.Request) {

}
