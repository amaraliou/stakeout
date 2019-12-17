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

// Student -> Struct to hold User information
type Student struct {
	gorm.Model
	User                  // Student Email to be verified (possibly use SheerID)
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	BirthDate      string `json:"birth_date" sql:"timestamp with time zone"`
	University     string `json:"university"`
	GraduationYear int    `json:"grad_year"`
}
