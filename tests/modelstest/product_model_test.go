package modelstest

import (
	"log"
	"testing"

	"gopkg.in/go-playground/assert.v1"
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
