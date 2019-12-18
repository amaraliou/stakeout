package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
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
	gorm.Model
	BasicStudent
	Degree         string `json:"degree"`
	GraduationYear int    `json:"grad_year"`
}
