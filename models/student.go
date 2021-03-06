package models

import (
	"errors"
	"log"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/nyaruka/phonenumbers"
	"golang.org/x/crypto/bcrypt"
)

// Token -> Struct to hold auth token information
type Token struct {
	UserID uint
	jwt.StandardClaims
}

// User -> Struct to hold basic user information
type User struct {
	Email      string `json:"email" gorm:"unique;not null"` // to add  gorm:"unique;not null"
	Password   string `json:"password"`
	IsVerified bool   `json:"verified"`
}

// Student -> struct to hold all the User information
type Student struct {
	Base
	User                  // Student Email to be verified (possibly use SheerID)
	IsStudent      bool   `json:"is_student"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	BirthDate      string `json:"birth_date" sql:"timestamp with time zone"`
	University     string `json:"university"`
	MobileNumber   string `json:"mobile_number"`
	CountryCode    string `json:"country"`
	GraduationYear int    `json:"grad_year"`
	Points         int    `json:"points"`
}

// Hash -> Generate hash for given password
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassword -> Verify a password given it's hash
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// BeforeSave will check hashes for passwords
func (student *Student) BeforeSave() error {
	hashedPassword, err := Hash(student.Password)
	if err != nil {
		return err
	}
	student.Password = string(hashedPassword)
	return nil
}

// Validate will validate the entries of the given student
func (student *Student) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if student.CountryCode != "" && student.MobileNumber != "" {
			if _, err := phonenumbers.Parse(student.MobileNumber, student.CountryCode); err != nil {
				return errors.New("Phone number ain't valid")
			}
		}

		return nil

	case "create":
		if student.Email == "" {
			return errors.New("Required Email")
		}

		if student.Password == "" {
			return errors.New("Required Password")
		}

		if err := checkmail.ValidateFormat(student.Email); err != nil {
			return errors.New("Invalid Email")
		}
		if student.CountryCode == "" {
			return errors.New("Required Country Code")
		}

		if student.MobileNumber == "" {
			return errors.New("Required Phone Number")
		}

		if _, err := phonenumbers.Parse(student.MobileNumber, student.CountryCode); err != nil {
			return errors.New("Phone number ain't valid")
		}

		return nil

	default:
		if student.Email == "" {
			return errors.New("Required Email")
		}

		if student.Password == "" {
			return errors.New("Required Password")
		}

		if err := checkmail.ValidateFormat(student.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil
	}
}

// CreateStudent -> Function to create a new student
func (student *Student) CreateStudent(db *gorm.DB) (*Student, error) {

	student.Points = 0

	err := db.Debug().Create(&student).Error
	if err != nil {
		return &Student{}, err
	}

	return student, nil
}

// FindAllStudents -> Function to retrieve all students
func (student *Student) FindAllStudents(db *gorm.DB) (*[]Student, error) {

	students := []Student{}
	err := db.Debug().Model(&Student{}).Limit(100).Find(&students).Error
	if err != nil {
		return &[]Student{}, err
	}

	return &students, nil
}

// FindStudentByID -> Function to retrieve a student given its ID
func (student *Student) FindStudentByID(db *gorm.DB, id string) (*Student, error) {

	err := db.Debug().Model(Student{}).Where("id = ?", id).Take(&student).Error
	if gorm.IsRecordNotFoundError(err) {
		return &Student{}, errors.New("Student not found")
	}

	if err != nil {
		return &Student{}, err
	}

	return student, nil
}

// UpdateStudent -> Function to update a given student
func (student *Student) UpdateStudent(db *gorm.DB, id string) (*Student, error) {

	err := student.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}

	err = db.Debug().Model(Student{}).Updates(&student).Error
	if err != nil {
		return &Student{}, err
	}

	return student.FindStudentByID(db, id)
}

// DeleteStudent -> Function to delete a student
func (student *Student) DeleteStudent(db *gorm.DB, id string) (int64, error) {

	db = db.Debug().Model(&Student{}).Where("id = ?", id).Take(&Student{}).Delete(&Student{})
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}
