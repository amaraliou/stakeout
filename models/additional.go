package models

import (
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
	City         string `json:"city"`
	County       string `json:"county"`
	Postcode     string `json:"postcode"` // To implement postcode validation
}
