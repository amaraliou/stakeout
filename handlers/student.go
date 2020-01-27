package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/amaraliou/stakeout/models"
	"github.com/amaraliou/stakeout/responses"
	"github.com/amaraliou/stakeout/auth"
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

	vars := mux.Vars(request)
	studentID := vars["id"]
	student := models.Student{}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
	}

	err = json.Unmarshal(body, &student)
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	tokenID, err := auth.ExtractTokenID(request)
	if err != nil {
		responses.ERROR(writer, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if tokenID != studentID {
		responses.ERROR(writer, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	err = student.Validate("update")
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	updatedStudent, err := student.UpdateStudent(server.DB, studentID)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(writer, http.StatusOK, updatedStudent)
}

// DeleteStudent -> handles DELETE /api/v1/student/<id:uuid>
func (server *Server) DeleteStudent(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	studentID := vars["id"]
	student := models.Student{}

	tokenID, err := auth.ExtractTokenID(request)
	if err != nil {
		responses.ERROR(writer, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if tokenID != "" && tokenID != studentID {
		responses.ERROR(writer, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	_, err = student.DeleteStudent(server.DB, studentID)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	writer.Header().Set("Entity", fmt.Sprintf("%s", studentID))
	responses.JSON(writer, http.StatusNoContent, "")
}
