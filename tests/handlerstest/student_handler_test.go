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

func TestCreateStudent(t *testing.T) {

	err := refreshStudentTable()
	if err != nil {
		log.Fatal(err)
	}

	samples := []struct {
		inputJSON    string
		statusCode   int
		nickname     string
		mobileNumber string
		email        string
		errorMessage string
	}{
		{
			inputJSON:    `{"email":"2310549a@student.gla.ac.uk", "password": "password", "country": "GB", "mobile_number":"07547775660"}`,
			statusCode:   201,
			email:        "2310549a@student.gla.ac.uk",
			mobileNumber: "07547775660",
			errorMessage: "",
		},
		{
			inputJSON:    `{"email":"", "password": "password", "country": "GB", "mobile_number":"07547775660"}`,
			statusCode:   422,
			errorMessage: "Required Email",
		},
		{
			inputJSON:    `{"email":"2310549a@student.gla.ac.uk", "password": "", "country": "GB", "mobile_number":"07547775660"}`,
			statusCode:   422,
			errorMessage: "Required Password",
		},
		{
			inputJSON:    `{"email":"2310549astudent.gla.ac.uk", "password": "password", "country": "GB", "mobile_number":"07547775660"}`,
			statusCode:   422,
			errorMessage: "Invalid Email",
		},
		{
			inputJSON:    `{"email":"2310549a@student.gla.ac.uk", "password": "password", "country": "", "mobile_number":"07547775660"}`,
			statusCode:   422,
			errorMessage: "Required Country Code",
		},
		{
			inputJSON:    `{"email":"2310549a@student.gla.ac.uk", "password": "password", "country": "GB", "mobile_number":""}`,
			statusCode:   422,
			errorMessage: "Required Phone Number",
		},
		{
			inputJSON:    `{"email":"2310549a@student.gla.ac.uk", "password": "password", "country": "GBRD", "mobile_number":"07560"}`,
			statusCode:   422,
			errorMessage: "Phone number ain't valid",
		},
		{
			inputJSON:    `{"email":"2310549a@student.gla.ac.uk", "password": "password", "country": "GB", "mobile_number":"07547775660"}`,
			statusCode:   500,
			errorMessage: "pq: duplicate key value violates unique constraint \"students_email_key\"",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/students", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateStudent)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["mobile_number"], v.mobileNumber)
			assert.Equal(t, responseMap["email"], v.email)
		}
		if v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetStudents(t *testing.T) {

	err := refreshStudentTable()
	if err != nil {
		log.Fatal(err)
	}

	_, err = seedStudents()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/students", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetStudents)
	handler.ServeHTTP(rr, req)

	var students []models.Student
	err = json.Unmarshal([]byte(rr.Body.String()), &students)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(students), 2)
}

func TestGetStudentByID(t *testing.T) {

	err := refreshStudentTable()
	if err != nil {
		log.Fatal(err)
	}

	student, err := seedOneStudent()
	if err != nil {
		log.Fatal(err)
	}

	studentSample := []struct {
		id         string
		statusCode int
		email      string
		firstName  string
		lastName   string
	}{
		{
			id:         student.ID.String(),
			statusCode: 200,
			email:      student.Email,
			firstName:  student.FirstName,
			lastName:   student.LastName,
		},
	}

	for _, v := range studentSample {

		req, err := http.NewRequest("GET", "/students", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}

		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetStudentByID)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, student.Email, responseMap["email"])
			assert.Equal(t, student.FirstName, responseMap["first_name"])
			assert.Equal(t, student.LastName, responseMap["last_name"])
		}
	}
}

func TestUpdateStudent(t *testing.T) {

	var AuthEmail, AuthPassword, AuthID string

	err := refreshStudentTable()
	if err != nil {
		log.Fatal(err)
	}

	students, err := seedStudents()
	if err != nil {
		log.Fatal(err)
	}

	authStudent := students[0]
	AuthID = authStudent.ID.String()
	AuthEmail = authStudent.Email
	AuthPassword = "password"
	unauthStudent := students[1]

	token, err := server.SignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		id           string
		updateJSON   string
		updateEmail  string
		updateNumber string
		statusCode   int
		tokenGiven   string
		errorMessage string
	}{
		{
			id:           AuthID,
			updateJSON:   `{"country": "GB", "mobile_number":"07564356660"}`,
			updateNumber: "07564356660",
			statusCode:   200,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			id:           AuthID,
			updateJSON:   `{"country": "GBDR", "mobile_number":"07564356660"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Phone number ain't valid",
		},
		{
			id:           unauthStudent.ID.String(),
			updateJSON:   `{"country": "GB", "mobile_number":"07564356660"}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
		// More cases to cover
	}

	for _, v := range samples {
		req, err := http.NewRequest("PUT", "/users", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}

		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateStudent)
		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["mobile_number"], v.updateNumber)
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
