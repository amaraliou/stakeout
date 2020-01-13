package handlerstest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amaraliou/apetitoso/models"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/assert.v1"
)

func refreshShopTable() error {
	err := server.DB.DropTableIfExists(&models.Shop{}).Error
	if err != nil {
		return err
	}

	err = server.DB.AutoMigrate(&models.Shop{}).Error
	if err != nil {
		return err
	}

	return nil
}

func seedOneShop() (models.Shop, error) {

	refreshShopTable()
	refreshAdminTable()

	admins, err := seedAdmins()
	if err != nil {
		log.Fatal(err)
	}

	shop := models.Shop{
		Name:        "Starbucks",
		Logo:        "https://starbucks.com",
		Description: "blah blah blah",
		Latitude:    45.5,
		Longitude:   45.5,
		ShopAddress: models.ShopAddress{
			Postcode:      "G12 8BG",
			AddressNumber: 10,
			AddressLine1:  "Bruh Street",
			TownOrCity:    "Bruh Town",
		},
	}

	admin := admins[0]

	err = server.DB.Model(&models.Shop{}).Create(&shop).Error
	if err != nil {
		return models.Shop{}, err
	}

	admin.ShopID = shop.ID

	err = server.DB.Model(models.Admin{}).Updates(&admin).Error
	if err != nil {
		return models.Shop{}, err
	}

	return shop, nil
}

func TestCreateShop(t *testing.T) {

	var AuthEmail, AuthPassword, AuthID string

	err := refreshShopTable()
	if err != nil {
		log.Fatal(err)
	}

	err = refreshAdminTable()
	if err != nil {
		log.Fatal(err)
	}

	admins, err := seedAdmins()
	if err != nil {
		log.Fatal(err)
	}

	authAdmin := admins[0]
	AuthID = authAdmin.ID.String()
	AuthEmail = authAdmin.Email
	AuthPassword = "password"

	token, err := server.AdminSignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		id           string
		createJSON   string
		statusCode   int
		tokenGiven   string
		shopName     string
		shopPostcode string
		errorMessage string
	}{
		{
			id:           AuthID,
			createJSON:   `{"name":"Some random shop", "description":"Random shop for testing", "postcode":"G12 8BY", "number":8, "address_1": "Amar Street", "town_or_city":"Glasgow"}`,
			statusCode:   201,
			tokenGiven:   tokenString,
			shopName:     "Some random shop",
			shopPostcode: "G12 8BY",
			errorMessage: "",
		},
		{
			id:           AuthID,
			createJSON:   `{"name":"", "description":"Random shop for testing", "postcode":"G12 8BY", "number":8, "address_1": "Amar Street", "town_or_city":"Glasgow"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required shop name",
		},
		{
			id:           AuthID,
			createJSON:   `{"name":"Some random shop", "description":"", "postcode":"G12 8BY", "number":8, "address_1": "Amar Street", "town_or_city":"Glasgow"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required shop description",
		},
		{
			id:           AuthID,
			createJSON:   `{"name":"Some random shop", "description":"bruh", "postcode":"", "number":8, "address_1": "Amar Street", "town_or_city":"Glasgow"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required shop postcode",
		},
		{
			id:           AuthID,
			createJSON:   `{"name":"Some random shop", "description":"bruh", "postcode":"G12 *BY", "number":0, "address_1": "Amar Street", "town_or_city":"Glasgow"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required shop address number",
		},
		{
			id:           AuthID,
			createJSON:   `{"name":"Some random shop", "description":"bruh", "postcode":"G12 *BY", "number":8, "address_1": "", "town_or_city":"Glasgow"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required shop address line 1",
		},
		{
			id:           AuthID,
			createJSON:   `{"name":"Some random shop", "description":"bruh", "postcode":"G12 *BY", "number":8, "address_1": "Amar Street", "town_or_city":""}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required town or city",
		},
		// More cases to cover
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/admins/", bytes.NewBufferString(v.createJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}

		req = mux.SetURLVars(req, map[string]string{"admin_id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateShop)
		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["name"], v.shopName)
			assert.Equal(t, responseMap["postcode"], v.shopPostcode)
		}

		if v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}

}
