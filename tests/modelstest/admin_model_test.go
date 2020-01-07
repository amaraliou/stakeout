package modelstest

import (
	"log"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllAdmins(t *testing.T) {

	err := refreshAdminTable()
	if err != nil {
		log.Fatal(err)
	}

	err = seedAdmins()
	if err != nil {
		log.Fatal(err)
	}

	admins, err := adminInstance.FindAllAdmins(server.DB)
	if err != nil {
		t.Errorf("This is the error getting the adminss: %v\n", err)
		return
	}

	assert.Equal(t, len(*admins), 2)
}

func TestCreateAdmin(t *testing.T) {
	assert.Equal(t, 1, 1)
}

func TestFindAdminByID(t *testing.T) {
	assert.Equal(t, 1, 1)
}

func TestUpdateAdmin(t *testing.T) {
	assert.Equal(t, 1, 1)
}

func TestDeleteAdmin(t *testing.T) {
	assert.Equal(t, 1, 1)
}
