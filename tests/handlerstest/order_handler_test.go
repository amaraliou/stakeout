package handlerstest

import (
	"testing"

	"github.com/amaraliou/stakeout/models"
	"gopkg.in/go-playground/assert.v1"
)

func refreshOrderTable() error {
	err := server.DB.DropTableIfExists(&models.Order{}).Error
	if err != nil {
		return err
	}

	err = server.DB.AutoMigrate(&models.Order{}).Error
	if err != nil {
		return err
	}

	return nil
}

func seedOneOrder() (models.Order, error) {

	refreshEverything()
	var total float32

	student, err := seedOneStudent()
	if err != nil {
		return models.Order{}, err
	}

	products, err := seedProducts()
	if err != nil {
		return models.Order{}, err
	}

	total = 0.0

	for _, product := range products {
		total = total + product.Price
	}

	order := models.Order{
		UserID:     student.ID,
		ShopID:     products[0].ShopID,
		OrderItems: products,
		OrderTotal: total,
	}

	err = server.DB.Model(&models.Order{}).Create(&order).Error
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func seedOrders() ([]models.Order, error) {

	refreshEverything()
	var total1 float32
	var total2 float32

	student, err := seedOneStudent()
	if err != nil {
		return []models.Order{}, err
	}

	products, err := seedProducts()
	if err != nil {
		return []models.Order{}, err
	}

	total1 = products[0].Price
	total2 = 0.0

	for _, product := range products {
		total2 = total2 + product.Price
	}

	orders := []models.Order{
		models.Order{
			UserID: student.ID,
			ShopID: products[0].ShopID,
			OrderItems: []models.Product{
				products[0],
			},
			OrderTotal: total1,
		},
		models.Order{
			UserID:     student.ID,
			ShopID:     products[0].ShopID,
			OrderItems: products,
			OrderTotal: total2,
		},
	}

	for _, order := range orders {
		err = server.DB.Model(&models.Order{}).Create(&order).Error
		if err != nil {
			return []models.Order{}, err
		}
	}

	return orders, nil
}

func TestCreateOrder(t *testing.T) {
	assert.Equal(t, 1, 1)
}

func TestGetOrders(t *testing.T) {
	assert.Equal(t, 1, 1)
}

func TestGetOrdersByShop(t *testing.T) {
	assert.Equal(t, 1, 1)
}

func TestGetOrdersByStudent(t *testing.T) {
	assert.Equal(t, 1, 1)
}

func TestGetOrderByID(t *testing.T) {
	assert.Equal(t, 1, 1)
}

func TestUpdateOrder(t *testing.T) {
	assert.Equal(t, 1, 1)
}

func TestDeleteOrder(t *testing.T) {
	assert.Equal(t, 1, 1)
}
