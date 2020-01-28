package models

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

const (
	OrderPending   uint8 = 0
	OrderPayed     uint8 = 1
	OrderReceived  uint8 = 2
	OrderConfirmed uint8 = 3
	OrderRefunding uint8 = 4
	OrderRefunded  uint8 = 5
	OrderCancel    uint8 = 6
)

var statusScope = []uint8{
	OrderPending,
	OrderPayed,
	OrderReceived,
	OrderConfirmed,
	OrderRefunding,
	OrderRefunded,
	OrderCancel,
}

// Order -> Struct to hold information about a specific order from a customer
type Order struct {
	Base
	UserID      uuid.UUID `json:"-" gorm:"user_id"`
	OrderedBy   Student   `json:"ordered_by" gorm:"foreignkey:UserID"`
	ShopID      uuid.UUID `json:"-" gorm:"shop_id"`
	OrderedFrom Shop      `json:"ordered_from" gorm:"foreignkey:ShopID"`
	OrderItems  []Product `json:"ordered_items" gorm:"many2many:order_products;"`
	OrderTotal  float32   `json:"total_price"`
	Status      uint8     `json:"status" gorm:"type:tinyint(1)"`
}

// Validate ...
func (order *Order) Validate(action string) error {
	switch strings.ToLower(action) {
	case "create":
		if order.UserID.String() == "00000000-0000-0000-0000-000000000000" {
			return errors.New("Required student")
		}

		if order.ShopID.String() == "00000000-0000-0000-0000-000000000000" {
			return errors.New("Required shop")
		}

	case "updatestatus":
		if order.Status < 0 {
			return errors.New("Invalid status")
		}

	default:
		return nil
	}

	return nil
}

// FindAllOrders ...
func (order *Order) FindAllOrders(db *gorm.DB) (*[]Order, error) {

	orders := []Order{}
	err := db.Debug().Model(&Order{}).Limit(100).Find(&orders).Error
	if err != nil {
		return &[]Order{}, err
	}

	return &orders, err
}

// FindAllOrdersByShop ...
func (order *Order) FindAllOrdersByShop(db *gorm.DB, shopID string) (*[]Order, error) {

	orders := []Order{}
	shop := Shop{}
	_, err := shop.FindShopByID(db, shopID)
	if err != nil {
		return &[]Order{}, err
	}

	err = db.Debug().Model(&Order{}).Where("shop_id = ?", shopID).Find(&orders).Error
	if err != nil {
		return &[]Order{}, err
	}

	return &orders, err
}

// FindAllOrdersByStudent ...
func (order *Order) FindAllOrdersByStudent(db *gorm.DB, studentID string) (*[]Order, error) {

	orders := []Order{}
	student := Student{}
	_, err := student.FindStudentByID(db, studentID)
	if err != nil {
		return &[]Order{}, err
	}

	err = db.Debug().Model(&Order{}).Where("user_id = ?", studentID).Find(&orders).Error
	if err != nil {
		return &[]Order{}, err
	}

	return &orders, err
}

// FindOrderByID ...
func (order *Order) FindOrderByID(db *gorm.DB, id string) (*Order, error) {

	err := db.Debug().Model(Order{}).Where("id = ?", id).Take(&order).Error
	if gorm.IsRecordNotFoundError(err) {
		return &Order{}, errors.New("Order not found")
	}

	if err != nil {
		return &Order{}, err
	}

	if order.UserID.String() != "00000000-0000-0000-0000-000000000000" {
		student := &Student{}
		err = db.Debug().Model(Student{}).Where("id = ?", order.UserID.String()).Take(&student).Error
		if err != nil {
			return order, errors.New("Student associated with this order not found")
		}
		order.OrderedBy = *student
	}

	if order.ShopID.String() != "00000000-0000-0000-0000-000000000000" {
		shop := &Shop{}
		err = db.Debug().Model(Shop{}).Where("id = ?", order.ShopID.String()).Take(&shop).Error
		if err != nil {
			return order, errors.New("Shop associated with this order not found")
		}
		order.OrderedFrom = *shop
	}

	return order, nil
}

// CreateOrder ...
func (order *Order) CreateOrder(db *gorm.DB) (*Order, error) {

	student := &Student{}
	err := db.Debug().Model(Student{}).Where("id = ?", order.UserID.String()).Take(&student).Error
	if err != nil {
		return &Order{}, errors.New("Student doesn't exist, can't create order")
	}

	shop := &Shop{}
	err = db.Debug().Model(Shop{}).Where("id = ?", order.ShopID.String()).Take(&shop).Error
	if err != nil {
		return &Order{}, errors.New("Shop doesn't exist, can't create order")
	}

	err = db.Debug().Create(&order).Error
	if err != nil {
		return &Order{}, err
	}

	order.OrderedBy = *student
	order.OrderedFrom = *shop

	return &Order{}, nil
}

// UpdateOrder ...
func (order *Order) UpdateOrder(db *gorm.DB, id string) (*Order, error) {

	err := db.Debug().Model(Order{}).Updates(&order).Error
	if err != nil {
		return &Order{}, nil
	}
	// More to cover

	return order, err
}

// DeleteOrder ...
func (order *Order) DeleteOrder(db *gorm.DB, id string) (int64, error) {

	db = db.Debug().Model(&Order{}).Where("id = ?", id).Take(&Order{}).Delete(&Order{})
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}
