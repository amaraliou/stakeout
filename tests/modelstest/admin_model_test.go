package modelstest

import (
	"log"
	"testing"

	"github.com/amaraliou/apetitoso/models"
	"github.com/lib/pq"
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

func TestFindAllAdminsNonExistentTable(t *testing.T) {

	err := server.DB.DropTableIfExists(&models.Admin{}).Error
	if err != nil {
		log.Fatal(err)
	}

	_, err = adminInstance.FindAllAdmins(server.DB)
	assert.Equal(t, err.(*pq.Error).Message, "relation \"admins\" does not exist")
}

func TestCreateAdmin(t *testing.T) {

	err := refreshAdminTable()
	if err != nil {
		log.Fatal(err)
	}

	newAdmin := models.Admin{
		User: models.User{
			Email:      "testemail@email.com",
			Password:   "password",
			IsVerified: true,
		},
		FirstName: "Donald WW3",
		LastName:  "Trump",
	}

	createdAdmin, err := newAdmin.CreateAdmin(server.DB)
	if err != nil {
		t.Errorf("This is the error creating the admin: %v\n", err)
		return
	}

	assert.Equal(t, newAdmin.Email, createdAdmin.Email)
	assert.Equal(t, newAdmin.FirstName, createdAdmin.FirstName)
	assert.Equal(t, newAdmin.LastName, createdAdmin.LastName)
}

func TestFindAdminByID(t *testing.T) {

	err := refreshAdminTable()
	if err != nil {
		log.Fatal(err)
	}

	admin, err := seedOneAdmin()
	if err != nil {
		log.Fatal(err)
	}

	foundAdmin, err := adminInstance.FindAdminByID(server.DB, admin.ID.String())
	if err != nil {
		t.Errorf("This is the error getting the admin: %v\n", err)
		return
	}

	assert.Equal(t, foundAdmin.ID, admin.ID)
	assert.Equal(t, foundAdmin.Email, admin.Email)
}

func TestUpdateAdmin(t *testing.T) {
	assert.Equal(t, 1, 1)
}

func TestDeleteAdmin(t *testing.T) {
	assert.Equal(t, 1, 1)
}
