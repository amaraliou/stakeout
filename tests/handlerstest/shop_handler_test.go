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

	err := server.DB.Model(&models.Shop{}).Create(&shop).Error
	if err != nil {
		return models.Shop{}, err
	}

	return shop, nil
}

func seedShops() ([]models.Shop, error) {

	refreshShopTable()

	shops := []models.Shop{
		models.Shop{
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
		},
		models.Shop{
			Name:        "Costa",
			Logo:        "https://costa.com",
			Description: "blah blah blah",
			Latitude:    55.5,
			Longitude:   55.5,
			ShopAddress: models.ShopAddress{
				Postcode:      "G12 8BY",
				AddressNumber: 10,
				AddressLine1:  "Bruh Avenue",
				TownOrCity:    "Bruh City",
			},
		},
	}

	for i := range shops {
		err := server.DB.Model(&models.Shop{}).Create(&shops[i]).Error
		if err != nil {
			return []models.Shop{}, err
		}
	}

	return shops, nil
}

func TestCreateShop(t *testing.T) {

	var AuthEmail, AuthPassword, AuthID string
	err := refreshAdminTable()
	if err != nil {
		log.Fatal(err)
	}

	err = refreshShopTable()
	if err != nil {
		log.Fatal(err)
	}

	err = refreshStudentTable()
	if err != nil {
		log.Fatal(err)
	}

	student, err := seedOneStudent()
	if err != nil {
		log.Fatal(err)
	}

	err = server.DB.Model(&models.Admin{}).AddForeignKey("shop_id", "shops(id)", "CASCADE", "CASCADE").Error
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
	unauthAdmin := admins[1]

	studentToken, err := server.SignIn(student.Email, "password")
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	studentTokenString := fmt.Sprintf("Bearer %v", studentToken)

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
		{
			id:           AuthID,
			createJSON:   `{"name":"Some random shop", "description":"bruh", "postcode":"G12 *BY", "number":8, "address_1": "Amar Street", "town_or_city":"Glasgow"}`,
			statusCode:   422,
			tokenGiven:   "shfdjsds",
			errorMessage: "token contains an invalid number of segments",
		},
		{
			id:           unauthAdmin.ID.String(),
			createJSON:   `{"name":"Some random shop", "description":"bruh", "postcode":"G12 *BY", "number":8, "address_1": "Amar Street", "town_or_city":"Glasgow"}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
		{
			id:           AuthID,
			createJSON:   `{"name":"Some random shop", "description":"bruh", "postcode":"G12 *BY", "number":8, "address_1": "Amar Street", "town_or_city":"Glasgow"}`,
			statusCode:   401,
			tokenGiven:   studentTokenString,
			errorMessage: "Unauthorized: This is not an admin token",
		},
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
			admin, err := authAdmin.FindAdminByID(server.DB, AuthID)
			if err != nil {
				log.Fatal(err)
			}
			assert.Equal(t, responseMap["name"], v.shopName)
			assert.Equal(t, responseMap["postcode"], v.shopPostcode)
			assert.Equal(t, responseMap["ID"], admin.ShopID.String())
			assert.Equal(t, responseMap["postcode"], admin.Shop.Postcode)
		}

		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetShops(t *testing.T) {

	err := refreshShopTable()
	if err != nil {
		log.Fatal(err)
	}

	shops, err := seedShops()
	if err != nil {
		log.Fatal(err)
	}

	err = refreshAdminTable()
	if err != nil {
		log.Fatal(err)
	}

	err = server.DB.Model(&models.Admin{}).AddForeignKey("shop_id", "shops(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		log.Fatal(err)
	}

	admins, err := seedAdmins()
	if err != nil {
		log.Fatal(err)
	}

	for i := range shops {
		currentAdmin := models.Admin{
			ShopID: shops[i].ID,
		}

		admin, err := currentAdmin.UpdateAdmin(server.DB, admins[i].ID.String())
		if err != nil {
			log.Fatal(err)
		}

		admins[i] = *admin
	}

	req, err := http.NewRequest("GET", "/shops", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetShops)
	handler.ServeHTTP(rr, req)

	var receivedShops []models.Shop
	err = json.Unmarshal([]byte(rr.Body.String()), &receivedShops)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(shops), 2)
}

func TestGetShopByID(t *testing.T) {

	err := refreshShopTable()
	if err != nil {
		log.Fatal(err)
	}

	shop, err := seedOneShop()
	if err != nil {
		log.Fatal(err)
	}

	shopSample := []struct {
		id           string
		statusCode   int
		name         string
		postcode     string
		errorMessage string
	}{
		{
			id:         shop.ID.String(),
			statusCode: 200,
			name:       shop.Name,
			postcode:   shop.Postcode,
		},
		{
			id:           "jdsfksjdfj",
			statusCode:   500,
			errorMessage: "pq: invalid input syntax for type uuid: \"jdsfksjdfj\"",
		},
		{
			id:           "1b56f03e-823c-4861-bee3-223c82e91c1f",
			statusCode:   500,
			errorMessage: "Shop not found",
		},
	}

	for _, v := range shopSample {

		req, err := http.NewRequest("GET", "/shops", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}

		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetShopByID)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, shop.ID.String(), responseMap["ID"])
			assert.Equal(t, shop.Name, responseMap["name"])
			assert.Equal(t, shop.Postcode, responseMap["postcode"])
		}

		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestUpdateShop(t *testing.T) {

	var AuthEmail, AuthPassword, AuthID string
	err := refreshShopTable()
	if err != nil {
		log.Fatal(err)
	}

	shops, err := seedShops()
	if err != nil {
		log.Fatal(err)
	}

	err = refreshStudentTable()
	if err != nil {
		log.Fatal(err)
	}

	student, err := seedOneStudent()
	if err != nil {
		log.Fatal(err)
	}

	err = refreshAdminTable()
	if err != nil {
		log.Fatal(err)
	}

	err = server.DB.Model(&models.Admin{}).AddForeignKey("shop_id", "shops(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		log.Fatal(err)
	}

	admins, err := seedAdmins()
	if err != nil {
		log.Fatal(err)
	}

	unauthAdmin := admins[0]
	unauthShop := shops[0]

	studentToken, err := server.SignIn(student.Email, "password")
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	studentTokenString := fmt.Sprintf("Bearer %v", studentToken)

	for i := range shops {
		currentAdmin := models.Admin{
			ShopID: shops[i].ID,
			User: models.User{
				Password: "password",
			},
		}

		admin, err := currentAdmin.UpdateAdmin(server.DB, admins[i].ID.String())
		if err != nil {
			log.Fatal(err)
		}

		admins[i] = *admin
	}

	authAdmin := admins[1]
	shop := shops[1]
	shopID := shop.ID.String()
	AuthID = authAdmin.ID.String()
	AuthEmail = authAdmin.Email
	AuthPassword = "password"

	token, err := server.AdminSignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		shopID       string
		adminID      string
		updateJSON   string
		updateName   string
		statusCode   int
		tokenGiven   string
		errorMessage string
	}{
		{
			adminID:      AuthID,
			shopID:       shopID,
			updateJSON:   `{"name": "Some random shop 2"}`,
			updateName:   "Some random shop 2",
			statusCode:   200,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			adminID:      unauthAdmin.ID.String(),
			shopID:       shopID,
			updateJSON:   `{"name": "Some random shop 2"}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
		{
			adminID:      AuthID,
			shopID:       shopID,
			updateJSON:   `{"name": "Some random shop 2"}`,
			statusCode:   401,
			tokenGiven:   studentTokenString,
			errorMessage: "Unauthorized: This is not an admin token",
		},
		{
			adminID:      AuthID,
			shopID:       shopID,
			updateJSON:   `{"name": "Some random shop 2"}`,
			statusCode:   422,
			tokenGiven:   "sbjadkjasjdahsdgjlkjasdjkai",
			errorMessage: "token contains an invalid number of segments",
		},
		{
			adminID:      AuthID,
			shopID:       shopID,
			updateJSON:   `{"name": "Some random shop 2"}`,
			statusCode:   422,
			tokenGiven:   "",
			errorMessage: "token contains an invalid number of segments",
		},
		{
			adminID:      AuthID,
			shopID:       unauthShop.ID.String(),
			updateJSON:   `{"name": "Some random shop 2"}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized: You are not the admin for this shop",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("PUT", "/shops", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}

		req = mux.SetURLVars(req, map[string]string{
			"admin_id": v.adminID,
			"shop_id":  v.shopID,
		})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateShop)
		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["name"], v.updateName)
		}

		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteShop(t *testing.T) {

	var AuthEmail, AuthPassword, AuthID string
	err := refreshShopTable()
	if err != nil {
		log.Fatal(err)
	}

	shops, err := seedShops()
	if err != nil {
		log.Fatal(err)
	}

	err = refreshStudentTable()
	if err != nil {
		log.Fatal(err)
	}

	student, err := seedOneStudent()
	if err != nil {
		log.Fatal(err)
	}

	err = refreshAdminTable()
	if err != nil {
		log.Fatal(err)
	}

	err = server.DB.Model(&models.Admin{}).AddForeignKey("shop_id", "shops(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		log.Fatal(err)
	}

	admins, err := seedAdmins()
	if err != nil {
		log.Fatal(err)
	}

	unauthAdmin := admins[0]
	unauthShop := shops[0]

	studentToken, err := server.SignIn(student.Email, "password")
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	studentTokenString := fmt.Sprintf("Bearer %v", studentToken)

	for i := range shops {
		currentAdmin := models.Admin{
			ShopID: shops[i].ID,
			User: models.User{
				Password: "password",
			},
		}

		admin, err := currentAdmin.UpdateAdmin(server.DB, admins[i].ID.String())
		if err != nil {
			log.Fatal(err)
		}

		admins[i] = *admin
	}

	authAdmin := admins[1]
	shop := shops[1]
	shopID := shop.ID.String()
	AuthID = authAdmin.ID.String()
	AuthEmail = authAdmin.Email
	AuthPassword = "password"

	token, err := server.AdminSignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		adminID      string
		shopID       string
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			adminID:      AuthID,
			shopID:       unauthShop.ID.String(),
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized: You are not the admin for this shop",
		},
		{
			adminID:      AuthID,
			shopID:       shopID,
			statusCode:   204,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			adminID:      AuthID,
			shopID:       shopID,
			statusCode:   422,
			tokenGiven:   "",
			errorMessage: "token contains an invalid number of segments",
		},
		{
			adminID:      AuthID,
			shopID:       shopID,
			statusCode:   422,
			tokenGiven:   "kjsklfsjsjdodw",
			errorMessage: "token contains an invalid number of segments",
		},
		{
			adminID:      AuthID,
			shopID:       shopID,
			statusCode:   401,
			tokenGiven:   studentTokenString,
			errorMessage: "Unauthorized: This is not an admin token",
		},
		{
			adminID:      unauthAdmin.ID.String(),
			shopID:       shopID,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("DELETE", "/shops", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}

		req = mux.SetURLVars(req, map[string]string{
			"admin_id": v.adminID,
			"shop_id":  v.shopID,
		})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteShop)
		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
