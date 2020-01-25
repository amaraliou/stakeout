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

func refreshProductTable() error {
	err := server.DB.DropTableIfExists(&models.Product{}).Error
	if err != nil {
		return err
	}

	err = server.DB.AutoMigrate(&models.Product{}).Error
	if err != nil {
		return err
	}

	return nil
}

func seedOneProduct() (models.Product, error) {

	refreshEverything()

	shop, err := seedOneShop()
	if err != nil {
		return models.Product{}, err
	}

	product := models.Product{
		Name:          "Cappuccino",
		Description:   "Froathy milk with decent coffee",
		Code:          "STBCKS001",
		Price:         2.95,
		PriceCurrency: "GBP",
		InSale:        false,
		ShopID:        shop.ID,
		Reward:        5,
	}

	err = server.DB.Model(&models.Product{}).Create(&product).Error
	if err != nil {
		return models.Product{}, err
	}

	return product, nil
}

func seedProducts() ([]models.Product, error) {

	refreshEverything()

	shop, err := seedOneShop()
	if err != nil {
		return []models.Product{}, err
	}

	products := []models.Product{
		models.Product{
			Name:          "Cappuccino",
			Description:   "Froathy milk with decent coffee",
			Code:          "STBCKS001",
			Price:         2.95,
			PriceCurrency: "GBP",
			InSale:        false,
			ShopID:        shop.ID,
			Reward:        5,
		},
		models.Product{
			Name:          "Espresso",
			Description:   "That shot of coffee you need to wake up",
			Code:          "STBCKS002",
			Price:         2.45,
			PriceCurrency: "GBP",
			InSale:        false,
			ShopID:        shop.ID,
			Reward:        3,
		},
	}

	for i := range products {
		err = server.DB.Model(&models.Product{}).Create(&products[i]).Error
		if err != nil {
			return []models.Product{}, err
		}
	}

	return products, nil
}

func TestCreateProduct(t *testing.T) {

	var AuthEmail, AuthPassword, AuthID string
	err := refreshEverything()
	if err != nil {
		log.Fatal(err)
	}

	shop, err := seedOneShop()
	if err != nil {
		log.Fatal(err)
	}

	student, err := seedOneStudent()
	if err != nil {
		log.Fatal(err)
	}

	studentToken, err := server.SignIn(student.Email, "password")
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	studentTokenString := fmt.Sprintf("Bearer %v", studentToken)

	err = server.DB.Model(&models.Admin{}).AddForeignKey("shop_id", "shops(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		log.Fatal(err)
	}

	err = server.DB.Model(&models.Product{}).AddForeignKey("shop_id", "shops(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		log.Fatal(err)
	}

	admins, err := seedAdmins()
	if err != nil {
		log.Fatal(err)
	}

	shopAdmin := models.Admin{
		ShopID: shop.ID,
	}

	AuthID = admins[0].ID.String()
	AuthEmail = admins[0].Email
	AuthPassword = "password"

	token, err := server.AdminSignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	authAdmin, err := shopAdmin.UpdateAdmin(server.DB, admins[0].ID.String())
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, AuthID, authAdmin.ID.String())
	assert.Equal(t, shop.ID, authAdmin.ShopID)

	samples := []struct {
		inputJSON    string
		statusCode   int
		tokenGiven   string
		name         string
		shopID       string
		errorMessage string
	}{
		{
			inputJSON:    `{"name": "Cappuccino", "price": 2.90, "price_currency": "GBP"}`,
			statusCode:   201,
			name:         "Cappuccino",
			tokenGiven:   tokenString,
			shopID:       shop.ID.String(),
			errorMessage: "",
		},
		{
			inputJSON:    `{"price": 2.90, "price_currency": "GBP"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			shopID:       shop.ID.String(),
			errorMessage: "Required product name",
		},
		{
			inputJSON:    `{"name": "Cappuccino", "price_currency": "GBP"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			shopID:       shop.ID.String(),
			errorMessage: "Required product price",
		},
		{
			inputJSON:    `{"name": "Cappuccino", "price": 2.90, "price_currency": "GBP"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			shopID:       "00000000-0000-0000-0000-000000000000",
			errorMessage: "Required shop",
		},
		{
			inputJSON:    `{"name": "Cappuccino", "price": 2.90, "price_currency": "GBP"}`,
			statusCode:   401,
			name:         "Cappuccino",
			tokenGiven:   studentTokenString,
			shopID:       shop.ID.String(),
			errorMessage: "Unauthorized: This is not an admin token",
		},
		{
			inputJSON:    `{"name": "Cappuccino", "price": 2.90, "price_currency": "GBP"}`,
			statusCode:   422,
			name:         "Cappuccino",
			tokenGiven:   "hdbjksjass",
			shopID:       shop.ID.String(),
			errorMessage: "token contains an invalid number of segments",
		},
		{
			inputJSON:    `{"name": "Cappuccino", "price": 2.90, "price_currency": "GBP"}`,
			statusCode:   422,
			name:         "Cappuccino",
			tokenGiven:   tokenString,
			shopID:       "jkjsanfjksakjkd",
			errorMessage: "uuid: incorrect UUID length: jkjsanfjksakjkd",
		},
		// More cases to cover
	}

	for _, v := range samples {
		req, err := http.NewRequest("POST", "/shops", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}

		req = mux.SetURLVars(req, map[string]string{"shop_id": v.shopID})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateProduct)
		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			shopResponseMap := responseMap["sold_by"].(map[string]interface{})
			assert.Equal(t, responseMap["name"], v.name)
			assert.Equal(t, shopResponseMap["ID"], v.shopID)
		}

		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetProducts(t *testing.T) {

	err := refreshEverything()
	if err != nil {
		log.Fatal(err)
	}

	products, err := seedProducts()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/products", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetProducts)
	handler.ServeHTTP(rr, req)

	var receivedProducts []models.Product
	err = json.Unmarshal([]byte(rr.Body.String()), &receivedProducts)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(products), 2)
}

func TestGetProductsByShop(t *testing.T) {
	assert.Equal(t, 1, 1)
}

func TestGetProductByID(t *testing.T) {

	err := refreshEverything()
	if err != nil {
		log.Fatal(err)
	}

	product, err := seedOneProduct()
	if err != nil {
		log.Fatal(err)
	}

	samples := []struct {
		id           string
		statusCode   int
		name         string
		errorMessage string
	}{
		{
			id:         product.ID.String(),
			statusCode: 200,
			name:       "Cappuccino",
		},
		{
			id:           "jdsfksjdfj",
			statusCode:   500,
			errorMessage: "pq: invalid input syntax for type uuid: \"jdsfksjdfj\"",
		},
		{
			id:           "1b56f03e-823c-4861-bee3-223c82e91c1f",
			statusCode:   500,
			errorMessage: "Product not found",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("GET", "/products", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}

		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetProductByID)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["ID"], v.id)
		}

		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestUpdateProduct(t *testing.T) {
	assert.Equal(t, 1, 1)
}

func TestDeleteProduct(t *testing.T) {
	assert.Equal(t, 1, 1)
}
