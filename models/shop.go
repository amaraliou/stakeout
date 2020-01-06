package models

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
)

// Shop -> Struct to hold shop information (try to figure out how to handle shop reg.)
type Shop struct {
	Base
	Name        string  `json:"name"`
	Logo        string  `json:"logo_link"` // Amazon S3
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Description string  `json:"description"`
	ShopAddress
}

// Validate ...
func (shop *Shop) Validate(action string) error {
	switch strings.ToLower(action) {
	case "create":
		if shop.Name == "" {
			return errors.New("Required shop name")
		}

		if shop.Description == "" {
			return errors.New("Required shop description")
		}

		if shop.Postcode == "" {
			return errors.New("Required shop postcode")
		}

		if shop.AddressNumber == 0 {
			return errors.New("Required shop address number")
		}

		if shop.AddressLine1 == "" {
			return errors.New("Required shop address line 1")
		}

		if shop.TownOrCity == "" {
			return errors.New("Required town or city")
		}

		return nil

	default:
		return nil
	}
}

// CreateShop ...
func (shop *Shop) CreateShop(db *gorm.DB) (*Shop, error) {

	err := db.Debug().Create(&shop).Error
	if err != nil {
		return &Shop{}, err
	}

	return shop, nil
}

// FindAllShops ...
func (shop *Shop) FindAllShops(db *gorm.DB) (*[]Shop, error) {

	shops := []Shop{}
	err := db.Debug().Model(&Shop{}).Limit(100).Find(&shops).Error
	if err != nil {
		return &[]Shop{}, err
	}

	return &shops, nil
}

// FindShopByID ...
func (shop *Shop) FindShopByID(db *gorm.DB, id string) (*Shop, error) {

	err := db.Debug().Model(Shop{}).Where("id = ?", id).Take(&shop).Error
	if gorm.IsRecordNotFoundError(err) {
		return &Shop{}, errors.New("Shop not found")
	}

	if err != nil {
		return &Shop{}, err
	}

	return shop, nil
}

// UpdateShop ...
func (shop *Shop) UpdateShop(db *gorm.DB, id string) (*Shop, error) {

	err := db.Debug().Model(Shop{}).Updates(&shop).Error
	if err != nil {
		return &Shop{}, err
	}

	return shop.FindShopByID(db, id)
}

// DeleteShop ...
func (shop *Shop) DeleteShop(db *gorm.DB, id string) (int64, error) {

	db = db.Debug().Model(&Shop{}).Where("id = ?", id).Take(&Shop{}).Delete(&Shop{})
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}
