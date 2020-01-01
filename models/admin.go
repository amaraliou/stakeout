package models

// Admin -> simple admin struct
type Admin struct {
	Base
	User
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
