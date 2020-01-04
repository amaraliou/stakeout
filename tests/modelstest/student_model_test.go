package modelstest

import (
	"log"
	"testing"

	"github.com/amaraliou/apetitoso/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllStudents(t *testing.T) {

	err := refreshStudentTable()
	if err != nil {
		log.Fatal(err)
	}

	err = seedStudents()
	if err != nil {
		log.Fatal(err)
	}

	students, err := studentInstance.FindAllStudents(server.DB)
	if err != nil {
		t.Errorf("This is the error getting the students: %v\n", err)
		return
	}

	assert.Equal(t, len(*students), 2)
}

func TestSaveStudent(t *testing.T) {

	err := refreshStudentTable()
	if err != nil {
		log.Fatal(err)
	}

	newStudent := models.Student{
		User: models.User{
			Email:      "testemail@email.com",
			Password:   "password",
			IsVerified: true,
		},
		IsStudent:      true,
		FirstName:      "Donald",
		LastName:       "Trump",
		BirthDate:      "09/09/1956",
		University:     "Temple University",
		MobileNumber:   "07547775660",
		CountryCode:    "GB",
		GraduationYear: 2021,
		Points:         0,
	}
	savedStudent, err := newStudent.CreateStudent(server.DB)
	if err != nil {
		t.Errorf("This is the error creating the student: %v\n", err)
		return
	}

	assert.Equal(t, newStudent.Email, savedStudent.Email)
	assert.Equal(t, newStudent.FirstName, savedStudent.FirstName)
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

	foundStudent, err := studentInstance.FindStudentByID(server.DB, student.ID.String())
	if err != nil {
		t.Errorf("This is the error getting the student: %v\n", err)
		return
	}

	assert.Equal(t, foundStudent.ID, student.ID)
	assert.Equal(t, foundStudent.Email, student.Email)
}
