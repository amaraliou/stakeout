package models

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Product -> Struct to hold product information
type Product struct {
	Base
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Code         string    `json:"code"`
	Price        float32   `json:"price"`
	InSale       bool      `json:"is_in_sale"`
	Discount     int       `json:"discount"`
	DiscountUnit string    `json:"discount_unit"`
	SoldBy       Shop      `json:"sold_by" gorm:"foreignkey:ShopID"`
	ShopID       uuid.UUID `json:"-" gorm:"shop_id"`
	Reward       int       `json:"reward"`
}

// Validate ...
func (product *Product) Validate(action string) error {
	switch strings.ToLower(action) {
	case "create":
		if product.Name == "" {
			return errors.New("Required product name")
		}

		if product.Price == 0.0 {
			return errors.New("Required product price")
		}

	default:
		return nil
	}

	return nil
}

// FindAllProducts ...
func (product *Product) FindAllProducts(db *gorm.DB) (*[]Product, error) {

	products := []Product{}
	err := db.Debug().Model(&Product{}).Limit(100).Find(&products).Error
	if err != nil {
		return &[]Product{}, err
	}

	return &products, nil
}

// FindAllProductsByShop ...
func (product *Product) FindAllProductsByShop(db *gorm.DB, shopID string) (*[]Product, error) {

	products := []Product{}
	err := db.Debug().Model(&Product{}).Where("shop_id = ?", shopID).Find(&products).Error
	if err != nil {
		return &[]Product{}, err
	}

	return &products, nil
}

// FindProductByID ...
func (product *Product) FindProductByID(db *gorm.DB, id string) (*Product, error) {
	return &Product{}, nil
}

// CreateProduct ...
func (product *Product) CreateProduct(db *gorm.DB, shopID string) (*Product, error) {
	return &Product{}, nil
}

// UpdateProduct ...
func (product *Product) UpdateProduct(db *gorm.DB, id string) (*Product, error) {
	return &Product{}, nil
}

// DeleteProduct ...
func (product *Product) DeleteProduct(db *gorm.DB, id string) (int64, error) {
	return 0, nil
}
