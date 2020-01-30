package modelstest

import (
	"log"
	"testing"

	"github.com/amaraliou/stakeout/models"
	"github.com/lib/pq"
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

func TestCreateOrder(t *testing.T) {

	err := refreshEverything()
	if err != nil {
		log.Fatal(err)
	}

	products, err := seedProducts()
	if err != nil {
		log.Fatal(err)
	}

	student, err := seedOneStudent()
	if err != nil {
		log.Fatal(err)
	}

	newOrder := models.Order{
		UserID: student.ID,
		ShopID: products[0].ShopID,
		OrderItems: []models.Product{
			products[0],
		},
		OrderTotal: products[0].Price,
		Status:     0,
	}

	savedOrder, err := newOrder.CreateOrder(server.DB)
	if err != nil {
		t.Errorf("This is the error creating the order: %v\n", err)
		return
	}

	assert.Equal(t, newOrder.ShopID, savedOrder.ShopID)
	assert.Equal(t, newOrder.OrderTotal, savedOrder.OrderTotal)
}

func TestUpdateOrder(t *testing.T) {

	err := refreshEverything()
	if err != nil {
		log.Fatal(err)
	}

	order, err := seedOneOrder()
	if err != nil {
		log.Fatal(err)
	}

	orderUpdate := models.Order{
		Status: 2,
	}

	updatedOrder, err := orderUpdate.UpdateOrder(server.DB, order.ID.String())
	if err != nil {
		t.Errorf("This is the error updating the order: %v\n", err)
		return
	}

	assert.Equal(t, updatedOrder.ID, orderUpdate.ID)
	assert.NotEqual(t, order.Status, updatedOrder.Status)
}

func TestDeleteOrder(t *testing.T) {

	err := refreshEverything()
	if err != nil {
		log.Fatal(err)
	}

	order, err := seedOneOrder()
	if err != nil {
		log.Fatal(err)
	}

	isDeleted, err := orderInstance.DeleteOrder(server.DB, order.ID.String())
	if err != nil {
		t.Errorf("This is the error deleting the order: %v\n", err)
		return
	}

	assert.Equal(t, isDeleted, int64(1))
}

func TestNonExistentOrderTable(t *testing.T) {

	err := refreshEverything()
	if err != nil {
		log.Fatal(err)
	}

	err = server.DB.DropTableIfExists(&models.Order{}).Error
	if err != nil {
		log.Fatal(err)
	}

	student, err := seedOneStudent()
	if err != nil {
		log.Fatal(err)
	}

	shop, err := seedOneShop()
	if err != nil {
		log.Fatal(err)
	}

	fakeOrder := models.Order{
		UserID: student.ID,
		ShopID: shop.ID,
	}

	_, err = orderInstance.FindAllOrders(server.DB)
	assert.Equal(t, err.(*pq.Error).Message, "relation \"orders\" does not exist")

	_, err = orderInstance.FindAllOrdersByStudent(server.DB, student.ID.String())
	assert.Equal(t, err.(*pq.Error).Message, "relation \"orders\" does not exist")

	_, err = orderInstance.FindAllOrdersByShop(server.DB, shop.ID.String())
	assert.Equal(t, err.(*pq.Error).Message, "relation \"orders\" does not exist")

	_, err = orderInstance.FindOrderByID(server.DB, shop.ID.String())
	assert.Equal(t, err.(*pq.Error).Message, "relation \"orders\" does not exist")

	_, err = fakeOrder.CreateOrder(server.DB)
	assert.Equal(t, err.(*pq.Error).Message, "relation \"orders\" does not exist")

	_, err = fakeOrder.UpdateOrder(server.DB, shop.ID.String())
	assert.Equal(t, err.(*pq.Error).Message, "relation \"orders\" does not exist")

	_, err = fakeOrder.DeleteOrder(server.DB, shop.ID.String())
	assert.Equal(t, err.(*pq.Error).Message, "relation \"orders\" does not exist")
}
