package models

import (
	"errors"
	"log"

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

	err := db.Debug().Create(&admin).Error
	if err != nil {
		return &Admin{}, err
	}

	return admin, nil
}

// FindAllAdmins -> Function to retrieve all admins
func (admin *Admin) FindAllAdmins(db *gorm.DB) (*[]Admin, error) {

	admins := []Admin{}
	err := db.Debug().Model(&Admin{}).Limit(100).Find(&admins).Error
	if err != nil {
		return &[]Admin{}, err
	}

	return &admins, nil
}

// FindAllAdminsWithShopID -> Function to retrieve all admins of a given shop
func (admin *Admin) FindAllAdminsWithShopID(db *gorm.DB, shopID string) (*[]Admin, error) {

	admins := []Admin{}
	err := db.Debug().Model(&Admin{}).Where("shop_id = ?", shopID).Find(&admins).Error
	if err != nil {
		return &[]Admin{}, err
	}

	return &admins, nil
}

// FindAdminByID -> Function to retrieve an admin given its ID
func (admin *Admin) FindAdminByID(db *gorm.DB, id string) (*Admin, error) {

	err := db.Debug().Model(Admin{}).Where("id = ?", id).Take(&admin).Error
	if err != nil {
		return &Admin{}, err
	}

	return admin, nil
}

// UpdateAdmin -> Function to update a given admin
func (admin *Admin) UpdateAdmin(db *gorm.DB, id string) (*Admin, error) {

	err := admin.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}

	err = db.Debug().Model(Admin{}).Updates(&admin).Error
	if err != nil {
		return &Admin{}, err
	}

	return admin.FindAdminByID(db, id)
}

// DeleteAdmin -> Function to delete an admin
func (admin *Admin) DeleteAdmin(db *gorm.DB, id string) (int64, error) {

	shopID := admin.ShopID.String()
	shop := Shop{}
	originalDB := db
	db = db.Debug().Model(&Admin{}).Where("id = ?", id).Take(&Admin{}).Delete(&Admin{})
	if db.Error != nil {
		return 0, db.Error
	}

	remainingAdmins, err := admin.FindAllAdminsWithShopID(db, shopID)
	if err != nil {
		return db.RowsAffected, errors.New("Couldn't retrieve remaining admins")
	}

	if len(*remainingAdmins) == 0 {

		shopRowsAffected, err := shop.DeleteShop(originalDB, shopID)
		if err != nil {
			return db.RowsAffected, errors.New("Couldn't delete the associated shop")
		}

		return db.RowsAffected + shopRowsAffected, nil
	}

	return db.RowsAffected, nil
}
