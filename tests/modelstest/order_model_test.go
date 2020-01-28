package modelstest

import (
	"log"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

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
