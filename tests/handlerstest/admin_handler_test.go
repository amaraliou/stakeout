package handlerstest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amaraliou/stakeout/models"
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

	var AuthEmail, AuthPassword, AuthID string

	err := refreshAdminTable()
	if err != nil {
		log.Fatal(err)
	}

	err = refreshStudentTable()
	if err != nil {
		log.Fatal(err)
	}

	admins, err := seedAdmins()
	if err != nil {
		log.Fatal(err)
	}

	student, err := seedOneStudent()
	if err != nil {
		log.Fatal(err)
	}

	authAdmin := admins[0]
	AuthID = authAdmin.ID.String()
	AuthEmail = authAdmin.Email
	AuthPassword = "password"
	unauthAdmin := admins[1]

	studentToken, err := server.SignIn(student.Email, "password")
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	studentTokenString := fmt.Sprintf("Bearer %v", studentToken)

	token, err := server.AdminSignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		id              string
		updateJSON      string
		updateFirstName string
		updateLastName  string
		statusCode      int
		tokenGiven      string
		errorMessage    string
	}{
		{
			id:              AuthID,
			updateJSON:      `{"first_name": "Aziz", "last_name": "Bruh"}`,
			updateFirstName: "Aziz",
			updateLastName:  "Bruh",
			statusCode:      200,
			tokenGiven:      tokenString,
			errorMessage:    "",
		},
		{
			id:           unauthAdmin.ID.String(),
			updateJSON:   `{"first_name": "Aziz", "last_name": "Bruh"}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
		{
			id:           AuthID,
			updateJSON:   `{"first_name": "Aziz", "last_name": "Bruh"}`,
			statusCode:   422,
			tokenGiven:   "kjaskjdkfjksssd",
			errorMessage: "token contains an invalid number of segments",
		},
		{
			id:           AuthID,
			updateJSON:   `{"first_name": "Aziz", "last_name": "Bruh"}`,
			statusCode:   401,
			tokenGiven:   studentTokenString,
			errorMessage: "Unauthorized: This is not an admin token",
		},
		{
			id:           AuthID,
			updateJSON:   `{"email": "email1@email.com", "password": "password", "first_name": "Aziz", "last_name": "Bruh"}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "pq: duplicate key value violates unique constraint \"admins_email_key\"",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("PUT", "/admins", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}

		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateAdmin)
		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["first_name"], v.updateFirstName)
			assert.Equal(t, responseMap["last_name"], v.updateLastName)
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteAdmin(t *testing.T) {

	var AuthEmail, AuthPassword, AuthID string

	err := refreshAdminTable()
	if err != nil {
		log.Fatal(err)
	}

	err = refreshStudentTable()
	if err != nil {
		log.Fatal(err)
	}

	admins, err := seedAdmins()
	if err != nil {
		log.Fatal(err)
	}

	student, err := seedOneStudent()
	if err != nil {
		log.Fatal(err)
	}

	authAdmin := admins[0]
	AuthID = authAdmin.ID.String()
	AuthEmail = authAdmin.Email
	AuthPassword = "password"
	unauthAdmin := admins[1]

	studentToken, err := server.SignIn(student.Email, "password")
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	studentTokenString := fmt.Sprintf("Bearer %v", studentToken)

	token, err := server.AdminSignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		id           string
		statusCode   int
		tokenGiven   string
		errorMessage string
	}{
		{
			id:           AuthID,
			statusCode:   204,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			id:           AuthID,
			statusCode:   422,
			tokenGiven:   "ljfalksjhfjsklhd",
			errorMessage: "token contains an invalid number of segments",
		},
		{
			id:           unauthAdmin.ID.String(),
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
		{
			id:           AuthID,
			statusCode:   401,
			tokenGiven:   studentTokenString,
			errorMessage: "Unauthorized: This is not an admin token",
		},
	}

	for _, v := range samples {
		req, err := http.NewRequest("DELETE", "/admins", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}

		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteAdmin)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
