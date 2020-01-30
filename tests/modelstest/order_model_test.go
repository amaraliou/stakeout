package modelstest

import (
	"log"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllOrders(t *testing.T) {

	err := refreshEverything()
	if err != nil {
		log.Fatal(err)
	}

	_, err = seedOrders()
	if err != nil {
		log.Fatal(err)
	}

	orders, err := orderInstance.FindAllOrders(server.DB)
	if err != nil {
		t.Errorf("This is the error getting the orders: %v\n", err)
		return
	}

	assert.Equal(t, len(*orders), 2)
}

func TestFindOrdersByShop(t *testing.T) {

	err := refreshEverything()
	if err != nil {
		log.Fatal(err)
	}

	orders, err := seedOrders()
	if err != nil {
		log.Fatal(err)
	}

	ordersByShop, err := orderInstance.FindAllOrdersByShop(server.DB, orders[0].ShopID.String())
	if err != nil {
		t.Errorf("This is the error getting the orders by shop: %v\n", err)
		return
	}

	assert.Equal(t, len(*ordersByShop), 2)
}

func TestFindAllOrdersByStudent(t *testing.T) {

	err := refreshEverything()
	if err != nil {
		log.Fatal(err)
	}

	orders, err := seedOrders()
	if err != nil {
		log.Fatal(err)
	}

	ordersByStudent, err := orderInstance.FindAllOrdersByStudent(server.DB, orders[0].UserID.String())
	if err != nil {
		t.Errorf("This is the error getting the orders by student: %v\n", err)
		return
	}

	testOrders := *ordersByStudent

	assert.Equal(t, len(*ordersByStudent), 2)
	assert.Equal(t, orders[0].ShopID, testOrders[0].ShopID)
}

func TestFindOrderByID(t *testing.T) {

	err := refreshEverything()
	if err != nil {
		log.Fatal(err)
	}

	order, err := seedOneOrder()
	if err != nil {
		log.Fatal(err)
	}

	foundOrder, err := orderInstance.FindOrderByID(server.DB, order.ID.String())
	if err != nil {
		t.Errorf("This is the error getting the order: %v\n", err)
		return
	}

	assert.Equal(t, foundOrder.ID, order.ID)
	assert.Equal(t, foundOrder.ShopID, order.ShopID)
}
