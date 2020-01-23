package handlerstest

import (
	"testing"

	"github.com/amaraliou/apetitoso/models"
	"gopkg.in/go-playground/assert.v1"
)

func refreshProductTable() error {
	err := server.DB.DropTableIfExists(&models.Product{}).Error
	if err != nil {
		return err
	}

	err = server.DB.AutoMigrate(&models.Product{}).Error
	if err != nil {
		return err
	}

	return nil
}

func seedOneProduct() (models.Product, error) {

	refreshEverything()

	shop, err := seedOneShop()
	if err != nil {
		return models.Product{}, err
	}

	product := models.Product{
		Name:          "Cappuccino",
		Description:   "Froathy milk with decent coffee",
		Code:          "STBCKS001",
		Price:         2.95,
		PriceCurrency: "GBP",
		InSale:        false,
		ShopID:        shop.ID,
		Reward:        5,
	}

	err = server.DB.Model(&models.Product{}).Create(&product).Error
	if err != nil {
		return models.Product{}, err
	}

	return product, nil
}

func seedProducts() ([]models.Product, error) {

	refreshEverything()

	shop, err := seedOneShop()
	if err != nil {
		return []models.Product{}, err
	}

	products := []models.Product{
		models.Product{
			Name:          "Cappuccino",
			Description:   "Froathy milk with decent coffee",
			Code:          "STBCKS001",
			Price:         2.95,
			PriceCurrency: "GBP",
			InSale:        false,
			ShopID:        shop.ID,
			Reward:        5,
		},
		models.Product{
			Name:          "Espresso",
			Description:   "That shot of coffee you need to wake up",
			Code:          "STBCKS002",
			Price:         2.45,
			PriceCurrency: "GBP",
			InSale:        false,
			ShopID:        shop.ID,
			Reward:        3,
		},
	}

	for i := range products {
		err = server.DB.Model(&models.Product{}).Create(&products[i]).Error
		if err != nil {
			return []models.Product{}, err
		}
	}

	return products, nil
}

func TestCreateProduct(t *testing.T) {
	assert.Equal(t, 1, 1)
}

func TestGetProducts(t *testing.T) {
	assert.Equal(t, 1, 1)
}

func TestGetProductsByShop(t *testing.T) {
	assert.Equal(t, 1, 1)
}

func TestGetProductByID(t *testing.T) {
	assert.Equal(t, 1, 1)
}

func TestUpdateProduct(t *testing.T) {
	assert.Equal(t, 1, 1)
}

func TestDeleteProduct(t *testing.T) {
	assert.Equal(t, 1, 1)
}
