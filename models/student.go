package models

import (
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// Token -> Struct to hold auth token information
type Token struct {
	UserID uint
	jwt.StandardClaims
}

// User -> Struct to hold basic user information
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// BasicStudent -> Struct to hold basic User information
type BasicStudent struct {
	User                // Student Email to be verified (possibly use SheerID)
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	BirthDate    string `json:"birth_date" sql:"timestamp with time zone"`
	University   string `json:"university"`
	MobileNumber int    `json:"mobile_number"`
}

// Student -> struct to hold all the User information
type Student struct {
	Base
	BasicStudent
	Degree         string     `json:"degree"`
	GraduationYear int        `json:"grad_year"`
	Address        *Address   `json:"address"`
	OtherAddresses []*Address `json:"other_addresses"`
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
