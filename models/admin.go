package models

import uuid "github.com/satori/go.uuid"

// Admin -> simple admin struct
type Admin struct {
	Base
	User
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Shop      Shop      `json:"shop" gorm:"foreignkey:ShopID"`
	ShopID    uuid.UUID `json:"-"`
}

// Validate ...
func (admin *Admin) Validate(action string) error {
	return nil
}
