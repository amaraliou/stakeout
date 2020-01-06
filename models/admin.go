package models

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Admin -> simple admin struct
type Admin struct {
	Base
	User
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Shop      Shop      `json:"shop" gorm:"foreignkey:ShopID"`
	ShopID    uuid.UUID `json:"-" gorm:"shop_id"`
}

// Validate ...
func (admin *Admin) Validate(action string) error {
	return nil
}

// BeforeSave will check hashes for passwords
func (admin *Admin) BeforeSave() error {
	hashedPassword, err := Hash(admin.Password)
	if err != nil {
		return err
	}
	admin.Password = string(hashedPassword)
	return nil
}

// CreateAdmin -> Function to create a new admin
func (admin *Admin) CreateAdmin(db *gorm.DB) (*Admin, error) {
	// To implement
	return &Admin{}, nil
}

// FindAllAdmins -> Function to retrieve all admins
func (admin *Admin) FindAllAdmins(db *gorm.DB) (*[]Admin, error) {
	// To implement
	return &[]Admin{}, nil
}

// FindAllAdminsWithShopID -> Function to retrieve all admins of a given shop
func (admin *Admin) FindAllAdminsWithShopID(db *gorm.DB, shopID string) (*[]Admin, error) {
	// To implement
	return &[]Admin{}, nil
}

// GetAdminByID -> Function to retrieve an admin given its ID
func (admin *Admin) GetAdminByID(db *gorm.DB, id string) (*Admin, error) {
	// To implement
	return &Admin{}, nil
}

// UpdateAdmin -> Function to update a given admin
func (admin *Admin) UpdateAdmin(db *gorm.DB, id string) (*Admin, error) {
	// To implement
	return &Admin{}, nil
}

// DeleteAdmin -> Function to delete an admin
func (admin *Admin) DeleteAdmin(db *gorm.DB, id string) (int64, error) {
	// To implement
	// Case 1: Delete admin but there are other remaining admins for shop
	// Case 2: Delete admin and there are no more admins for shop -> delete shop
	return 0, nil
}
