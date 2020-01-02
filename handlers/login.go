package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/amaraliou/apetitoso/auth"
	"github.com/amaraliou/apetitoso/models"
	"github.com/amaraliou/apetitoso/responses"
	"golang.org/x/crypto/bcrypt"
)

// Login -> handles POST /api/v1/login
func (server *Server) Login(writer http.ResponseWriter, request *http.Request) {

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	student := models.Student{}
	err = json.Unmarshal(body, &student)
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	err = student.Validate("")
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	token, err := server.SignIn(student.Email, student.Password)
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	responses.JSON(writer, http.StatusOK, token)
}

// SignIn -> retrieves JWT token given username and password
func (server *Server) SignIn(email, password string) (string, error) {

	var err error
	student := models.Student{}
	err = server.DB.Debug().Model(models.Student{}).Where("email = ?", email).Take(&student).Error
	if err != nil {
		return "", err
	}

	err = models.VerifyPassword(student.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	return auth.CreateToken(student.ID)
}
