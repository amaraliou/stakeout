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
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Code          string    `json:"code"`
	Price         float32   `json:"price"`
	PriceCurrency string    `json:"price_currency"`
	InSale        bool      `json:"is_in_sale"`
	Discount      int       `json:"discount"`
	DiscountUnit  string    `json:"discount_unit"`
	SoldBy        Shop      `json:"sold_by" gorm:"foreignkey:ShopID"`
	ShopID        uuid.UUID `json:"-" gorm:"shop_id"`
	Reward        int       `json:"reward"`
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

		if product.ShopID.String() == "00000000-0000-0000-0000-000000000000" {
			return errors.New("Required shop")
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

	err := db.Debug().Model(Product{}).Where("id = ?", id).Take(&product).Error
	if gorm.IsRecordNotFoundError(err) {
		return &Product{}, errors.New("Product not found")
	}

	if err != nil {
		return &Product{}, err
	}

	if product.ShopID.String() != "00000000-0000-0000-0000-000000000000" {
		shop := &Shop{}
		err = db.Debug().Model(Shop{}).Where("id = ?", product.ShopID.String()).Take(&shop).Error
		if err != nil {
			return product, errors.New("Shop associated with this admin not found")
		}
		product.SoldBy = *shop
	}

	return product, nil
}

// CreateProduct ...
func (product *Product) CreateProduct(db *gorm.DB) (*Product, error) {

	shop := &Shop{}
	err := db.Debug().Model(Shop{}).Where("id = ?", product.ShopID.String()).Take(&shop).Error
	if err != nil {
		return &Product{}, errors.New("Shop doesn't exist, can't create product")
	}

	err = db.Debug().Create(&product).Error
	if err != nil {
		return &Product{}, err
	}

	product.SoldBy = *shop

	return product, nil
}

// UpdateProduct ...
func (product *Product) UpdateProduct(db *gorm.DB, id string) (*Product, error) {

	err := db.Debug().Model(Product{}).Updates(&product).Error
	if err != nil {
		return &Product{}, err
	}
	// More to cover

	return product, nil
}

// DeleteProduct ...
func (product *Product) DeleteProduct(db *gorm.DB, id string) (int64, error) {

	db = db.Debug().Model(&Product{}).Where("id = ?", id).Take(&Product{}).Delete(&Product{})
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}
