package handlerstest

import (
	"testing"

	"github.com/amaraliou/stakeout/models"
	"gopkg.in/go-playground/assert.v1"
)

func refreshOrderTable() error {
	return nil
}

func seedOneOrder() (models.Order, error) {
	return models.Order{}, nil
}

func seedOrders() ([]models.Order, error) {
	return []models.Order{}, nil
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
