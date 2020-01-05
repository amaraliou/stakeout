package modelstest

import (
	"errors"
	"log"
	"testing"

	"github.com/amaraliou/apetitoso/models"
	"github.com/lib/pq"
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

func TestFindAllStudentsNonExistingTable(t *testing.T) {
	err := server.DB.DropTableIfExists(&models.Student{}).Error
	if err != nil {
		log.Fatal(err)
	}

	_, err = studentInstance.FindAllStudents(server.DB)
	assert.Equal(t, err.(*pq.Error).Message, "relation \"students\" does not exist")
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

func TestUpdateStudent(t *testing.T) {

	err := refreshStudentTable()
	if err != nil {
		log.Fatal(err)
	}

	student, err := seedOneStudent()
	if err != nil {
		log.Fatal(err)
	}

	studentUpdate := models.Student{
		User: models.User{
			Email:      "email@email.com",
			Password:   "password",
			IsVerified: true,
		},
		IsStudent:      true,
		FirstName:      "Trump",
		LastName:       "Bruh",
		BirthDate:      "09/09/1956",
		University:     "Temple University",
		MobileNumber:   "07547775660",
		CountryCode:    "GB",
		GraduationYear: 2021,
		Points:         0,
	}

	updatedStudent, err := studentUpdate.UpdateStudent(server.DB, student.ID.String())
	if err != nil {
		t.Errorf("This is the error updating the student: %v\n", err)
		return
	}

	assert.Equal(t, updatedStudent.ID, studentUpdate.ID)
	assert.Equal(t, updatedStudent.FirstName, studentUpdate.FirstName)
}

func TestDeleteStudent(t *testing.T) {

	err := refreshStudentTable()
	if err != nil {
		log.Fatal(err)
	}

	student, err := seedOneStudent()
	if err != nil {
		log.Fatal(err)
	}

	isDeleted, err := studentInstance.DeleteStudent(server.DB, student.ID.String())
	if err != nil {
		t.Errorf("This is the error deleting the student: %v\n", err)
		return
	}

	assert.Equal(t, isDeleted, int64(1))
}

func TestDeleteWrongStudent(t *testing.T) {
	err := refreshStudentTable()
	if err != nil {
		log.Fatal(err)
	}

	_, err = seedOneStudent()
	if err != nil {
		log.Fatal(err)
	}

	_, err = studentInstance.DeleteStudent(server.DB, "8258e9fd-7769-4eb5-8b82-5f597e94e7a1")
	assert.Equal(t, err, errors.New("record not found"))
}
