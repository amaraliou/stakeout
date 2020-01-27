package modelstest

import (
	"github.com/amaraliou/stakeout/models"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"testing"
)

func TestFindAllProducts(t *testing.T) {

	err := refreshEverything()
	if err != nil {
		log.Fatal(err)
	}

	_, err = seedProducts()
	if err != nil {
		log.Fatal(err)
	}

	products, err := productInstance.FindAllProducts(server.DB)
	if err != nil {
		t.Errorf("This is the error getting the products: %v\n", err)
		return
	}

	assert.Equal(t, len(*products), 2)
}

func TestFindAllProductsByShop(t *testing.T) {

	err := refreshEverything()
	if err != nil {
		log.Fatal(err)
	}

	ps, err := seedProducts()
	if err != nil {
		log.Fatal(err)
	}

	shopID := ps[0].ShopID.String()

	products, err := productInstance.FindAllProductsByShop(server.DB, shopID)
	if err != nil {
		t.Errorf("This is the error getting the products by shop: %v\n", err)
		return
	}

	assert.Equal(t, len(*products), 2)
}

func TestFindProductByID(t *testing.T) {

	err := refreshEverything()
	if err != nil {
		log.Fatal(err)
	}

	product, err := seedOneProduct()
	if err != nil {
		log.Fatal(err)
	}

	foundProduct, err := productInstance.FindProductByID(server.DB, product.ID.String())
	if err != nil {
		t.Errorf("This is the error getting the product: %v\n", err)
		return
	}

	assert.Equal(t, foundProduct.ID, product.ID)
	assert.Equal(t, foundProduct.Name, product.Name)
	assert.Equal(t, foundProduct.Price, product.Price)
}

func TestCreateProduct(t *testing.T) {

	err := refreshEverything()
	if err != nil {
		log.Fatal(err)
	}

	shop, err := seedOneShop()
	if err != nil {
		log.Fatal(err)
	}

	newProduct := models.Product{
		Name:          "Cappuccino",
		Description:   "Froathy milk with decent coffee",
		Code:          "STBCKS001",
		Price:         2.95,
		PriceCurrency: "GBP",
		InSale:        false,
		ShopID:        shop.ID,
		Reward:        5,
	}

	savedProduct, err := newProduct.CreateProduct(server.DB)
	if err != nil {
		t.Errorf("This is the error creating the product: %v\n", err)
		return
	}

	assert.Equal(t, newProduct.Name, savedProduct.Name)
	assert.Equal(t, newProduct.ShopID, savedProduct.ShopID)
}

func TestUpdateProduct(t *testing.T) {

	err := refreshEverything()
	if err != nil {
		log.Fatal(err)
	}

	product, err := seedOneProduct()
	if err != nil {
		log.Fatal(err)
	}

	productUpdate := models.Product{
		Name:          "Macchiato",
		Description:   "No description, I can barely describe a macchiato",
		Code:          "STBCKS003",
		Price:         2.65,
		PriceCurrency: "GBP",
		InSale:        false,
		Reward:        5,
	}

	updatedProduct, err := productUpdate.UpdateProduct(server.DB, product.ID.String())
	if err != nil {
		t.Errorf("This is the error updating the product: %v\n", err)
		return
	}

	assert.Equal(t, updatedProduct.ID, productUpdate.ID)
	assert.Equal(t, updatedProduct.Name, productUpdate.Name)
}

func TestDeleteProduct(t *testing.T) {

	err := refreshEverything()
	if err != nil {
		log.Fatal(err)
	}

	product, err := seedOneProduct()
	if err != nil {
		log.Fatal(err)
	}

	isDeleted, err := productInstance.DeleteProduct(server.DB, product.ID.String())
	if err != nil {
		t.Errorf("This is the error deleting the product: %v\n", err)
		return
	}

	assert.Equal(t, isDeleted, int64(1))
}
