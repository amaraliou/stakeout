package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/amaraliou/apetitoso/models"
	"github.com/amaraliou/apetitoso/responses"
)

// CreateStudent -> handles POST /api/v1/student/
func (server *Server) CreateStudent(writer http.ResponseWriter, request *http.Request) {

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
	}

	student := models.Student{}
	err = json.Unmarshal(body, &student)
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	err = student.Validate("create")
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	studentCreated, err := student.CreateStudent(server.DB)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	writer.Header().Set("Location", fmt.Sprintf("%s%s/%s", request.Host, request.RequestURI, studentCreated.ID.String()))
	responses.JSON(writer, http.StatusCreated, studentCreated)
}

// GetStudents -> handles GET /api/v1/student/
func (server *Server) GetStudents(writer http.ResponseWriter, request *http.Request) {

	student := models.Student{}
	students, err := student.FindAllStudents(server.DB)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(writer, http.StatusOK, students)
}

// GetStudentByID -> handles GET /api/v1/student/<id:uuid>
func (server *Server) GetStudentByID(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	student := models.Student{}
	studentRetrieved, err := student.FindStudentByID(server.DB, vars["id"])
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(writer, http.StatusOK, studentRetrieved)
}

// UpdateStudent -> handles PUT /api/v1/student/<id:uuid>
func (server *Server) UpdateStudent(writer http.ResponseWriter, request *http.Request) {

}

// DeleteStudent -> handles DELETE /api/v1/student/<id:uuid>
func (server *Server) DeleteStudent(writer http.ResponseWriter, request *http.Request) {

}
