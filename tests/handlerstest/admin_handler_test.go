package handlerstest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

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

}

func TestGetAdminByID(t *testing.T) {

}

func TestUpdateAdmin(t *testing.T) {

}

func TestDeleteAdmin(t *testing.T) {

}
