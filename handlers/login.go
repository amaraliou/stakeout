package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/amaraliou/stakeout/models"
	"github.com/amaraliou/stakeout/responses"
	"github.com/amaraliou/stakeout/auth"
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

// AdminLogin -> handles POST /api/v1/admin/login
func (server *Server) AdminLogin(writer http.ResponseWriter, request *http.Request) {

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	admin := models.Admin{}
	err = json.Unmarshal(body, &admin)
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	err = admin.Validate("")
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	token, err := server.AdminSignIn(admin.Email, admin.Password)
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	responses.JSON(writer, http.StatusOK, token)
}

// SignIn -> retrieves user JWT token given username and password
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

// AdminSignIn -> retrieves admin JWT token given username and password
func (server *Server) AdminSignIn(email, password string) (string, error) {

	var err error
	admin := models.Admin{}
	err = server.DB.Debug().Model(&models.Admin{}).Where("email = ?", email).Take(&admin).Error
	if err != nil {
		return "", err
	}

	err = models.VerifyPassword(admin.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	return auth.CreateAdminToken(admin.ID)
}
