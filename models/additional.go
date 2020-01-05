package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// CurrentLocation -> struct to hold current location of a student
type CurrentLocation struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

// Address -> struct to hold address information
type Address struct {
	gorm.Model
	AddressLine1 string `json:"address_1"`
	AddressLine2 string `json:"address_2"`
	TownOrCity   string `json:"town_or_city"`
	County       string `json:"county"`
	Postcode     string `json:"postcode"` // To implement postcode validation
	Primary      bool   `json:"is_primary_address"`
}

// Validate address -> to configure
func (address *Address) Validate(action string) error {
	return nil
}

// CreateAddress -> Function to create a new address
func (address *Address) CreateAddress(db *gorm.DB) (*Address, error) {

	err := db.Debug().Create(&address).Error
	if err != nil {
		return &Address{}, err
	}

	return address, nil
}

// GetAddressByID -> Function to retrieve an address given its ID
func (address *Address) GetAddressByID(db *gorm.DB, id uint) (*Address, error) {

	err := db.Debug().Model(Address{}).Where("id = ?", id).Take(&address).Error
	if gorm.IsRecordNotFoundError(err) {
		return &Address{}, errors.New("Student not found")
	}

	if err != nil {
		return &Address{}, err
	}

	return address, nil
}
