package handlerstest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amaraliou/apetitoso/models"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateAdmin(t *testing.T) {

	err := refreshAdminTable()
	if err != nil {
		log.Fatal(err)
	}

	samples := []struct {
		inputJSON    string
		statusCode   int
		email        string
		firstName    string
		lastName     string
		errorMessage string
	}{
		{
			inputJSON:  `{"email":"admin@gmail.com", "password": "password", "first_name": "John", "last_name":"Doe"}`,
			statusCode: 201,
			email:      "admin@gmail.com",
			firstName:  "John",
			lastName:   "Doe",
		},
		{
			inputJSON:    `{"email":"", "password": "password", "first_name": "John", "last_name":"Doe"}`,
			statusCode:   422,
			errorMessage: "Required Email",
		},
		{
			inputJSON:    `{"email":"admin@gmail.com", "password": "", "first_name": "John", "last_name":"Doe"}`,
			statusCode:   422,
			errorMessage: "Required Password",
		},
		{
			inputJSON:    `{"email":"admingmail.com", "password": "password", "first_name": "John", "last_name":"Doe"}`,
			statusCode:   422,
			errorMessage: "Invalid Email",
		},
		{
			inputJSON:    `{"email":"admin@gmail.com", "password": "password", "first_name": "John", "last_name":"Doe"}`,
			statusCode:   500,
			errorMessage: "pq: duplicate key value violates unique constraint \"admins_email_key\"",
		},
	}

	for _, v := range samples {
		req, err := http.NewRequest("POST", "/admins", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateAdmin)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["first_name"], v.firstName)
			assert.Equal(t, responseMap["email"], v.email)
		}

		if v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetAdmins(t *testing.T) {

	err := refreshAdminTable()
	if err != nil {
		log.Fatal(err)
	}

	_, err = seedAdmins()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/admins", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetAdmins)
	handler.ServeHTTP(rr, req)

	var admins []models.Admin
	err = json.Unmarshal([]byte(rr.Body.String()), &admins)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(admins), 2)
}

func TestGetAdminByID(t *testing.T) {

	err := refreshAdminTable()
	if err != nil {
		log.Fatal(err)
	}

	admin, err := seedOneAdmin()
	if err != nil {
		log.Fatal(err)
	}

	adminSample := []struct {
		id           string
		statusCode   int
		email        string
		firstName    string
		lastName     string
		errorMessage string
	}{
		{
			id:         admin.ID.String(),
			statusCode: 200,
			email:      admin.Email,
			firstName:  admin.FirstName,
			lastName:   admin.LastName,
		},
		{
			id:           "jdsfksjdfj",
			statusCode:   500,
			errorMessage: "pq: invalid input syntax for type uuid: \"jdsfksjdfj\"",
		},
		{
			id:           "1b56f03e-823c-4861-bee3-223c82e91c1f",
			statusCode:   500,
			errorMessage: "Admin not found",
		},
	}

	for _, v := range adminSample {

		req, err := http.NewRequest("GET", "/admins", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}

		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetAdminByID)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, admin.Email, responseMap["email"])
			assert.Equal(t, admin.FirstName, responseMap["first_name"])
			assert.Equal(t, admin.LastName, responseMap["last_name"])
		}

		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestUpdateAdmin(t *testing.T) {

}

func TestDeleteAdmin(t *testing.T) {

}
