package modelstest

import (
	"log"
	"testing"

	"github.com/amaraliou/apetitoso/models"
	"github.com/lib/pq"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllShops(t *testing.T) {

	err := refreshShopTable()
	if err != nil {
		log.Fatal(err)
	}

	err = seedShops()
	if err != nil {
		log.Fatal(err)
	}

	shops, err := shopInstance.FindAllShops(server.DB)
	if err != nil {
		t.Errorf("This is the error getting the shops: %v\n", err)
		return
	}

	assert.Equal(t, len(*shops), 2)
}

func TestFindShopByID(t *testing.T) {

	err := refreshShopTable()
	if err != nil {
		log.Fatal(err)
	}

	shop, err := seedOneShop()
	if err != nil {
		log.Fatal(err)
	}

	foundShop, err := shopInstance.FindShopByID(server.DB, shop.ID.String())
	if err != nil {
		t.Errorf("This is the error getting the shop: %v\n", err)
		return
	}

	assert.Equal(t, foundShop.ID, shop.ID)
	assert.Equal(t, foundShop.Name, shop.Name)
	assert.Equal(t, foundShop.Postcode, shop.Postcode)
}

func TestCreateShop(t *testing.T) {

	err := refreshShopTable()
	if err != nil {
		log.Fatal(err)
	}

	newShop := models.Shop{
		Name:        "Costa",
		Logo:        "logo_link",
		Description: "blah blah blah",
		Latitude:    45.5,
		Longitude:   45.5,
		ShopAddress: models.ShopAddress{
			Postcode:      "G12 8BY",
			AddressNumber: 10,
			AddressLine1:  "Costa Street",
			TownOrCity:    "Costa Town",
		},
	}

	savedShop, err := newShop.CreateShop(server.DB)
	if err != nil {
		t.Errorf("This is the error creating the shop: %v\n", err)
		return
	}

	assert.Equal(t, newShop.Name, savedShop.Name)
	assert.Equal(t, newShop.Postcode, savedShop.Postcode)
}

func TestUpdateShop(t *testing.T) {

	err := refreshShopTable()
	if err != nil {
		log.Fatal(err)
	}

	shop, err := seedOneShop()
	if err != nil {
		log.Fatal(err)
	}

	shopUpdate := models.Shop{
		Name:        "Costa",
		Logo:        "logo_link",
		Description: "blah blah blah",
		Latitude:    45.5,
		Longitude:   45.5,
		ShopAddress: models.ShopAddress{
			Postcode:      "G12 8BY",
			AddressNumber: 10,
			AddressLine1:  "Costa Street",
			TownOrCity:    "Costa Town",
		},
	}

	updatedShop, err := shopUpdate.UpdateShop(server.DB, shop.ID.String())
	if err != nil {
		t.Errorf("This is the error updating the shop: %v\n", err)
		return
	}

	assert.Equal(t, updatedShop.ID, shopUpdate.ID)
	assert.Equal(t, updatedShop.Name, shopUpdate.Name)
}

func TestDeleteShop(t *testing.T) {

	err := refreshShopTable()
	if err != nil {
		log.Fatal(err)
	}

	shop, err := seedOneShop()
	if err != nil {
		log.Fatal(err)
	}

	isDeleted, err := shopInstance.DeleteShop(server.DB, shop.ID.String())
	if err != nil {
		t.Errorf("This is the error deleting the shop: %v\n", err)
		return
	}

	assert.Equal(t, isDeleted, int64(1))
}

func TestNonExistentShopsTable(t *testing.T) {

	err := server.DB.DropTableIfExists(&models.Shop{}).Error
	if err != nil {
		log.Fatal(err)
	}

	_, err = shopInstance.FindAllShops(server.DB)
	assert.Equal(t, err.(*pq.Error).Message, "relation \"shops\" does not exist")

	_, err = shopInstance.FindShopByID(server.DB, "random_id")
	assert.Equal(t, err.(*pq.Error).Message, "relation \"shops\" does not exist")

	_, err = shopInstance.CreateShop(server.DB)
	assert.Equal(t, err.(*pq.Error).Message, "relation \"shops\" does not exist")

	_, err = shopInstance.UpdateShop(server.DB, "random_id")
	assert.Equal(t, err.(*pq.Error).Message, "relation \"shops\" does not exist")

	_, err = shopInstance.DeleteShop(server.DB, "random_id")
	assert.Equal(t, err.(*pq.Error).Message, "relation \"shops\" does not exist")
}
