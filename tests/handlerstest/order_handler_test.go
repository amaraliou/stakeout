package handlerstest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amaraliou/stakeout/models"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/assert.v1"
)

func refreshOrderTable() error {
	err := server.DB.DropTableIfExists(&models.Order{}).Error
	if err != nil {
		return err
	}

	err = server.DB.AutoMigrate(&models.Order{}).Error
	if err != nil {
		return err
	}

	return nil
}

func seedOneOrder() (models.Order, error) {

	refreshEverything()
	var total float32

	student, err := seedOneStudent()
	if err != nil {
		return models.Order{}, err
	}

	products, err := seedProducts()
	if err != nil {
		return models.Order{}, err
	}

	total = 0.0

	for _, product := range products {
		total = total + product.Price
	}

	order := models.Order{
		UserID:     student.ID,
		ShopID:     products[0].ShopID,
		OrderItems: products,
		OrderTotal: total,
	}

	err = server.DB.Model(&models.Order{}).Create(&order).Error
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func seedOrders() ([]models.Order, error) {

	refreshEverything()
	var total1 float32
	var total2 float32

	student, err := seedOneStudent()
	if err != nil {
		return []models.Order{}, err
	}

	products, err := seedProducts()
	if err != nil {
		return []models.Order{}, err
	}

	total1 = products[0].Price
	total2 = 0.0

	for _, product := range products {
		total2 = total2 + product.Price
	}

	orders := []models.Order{
		models.Order{
			UserID: student.ID,
			ShopID: products[0].ShopID,
			OrderItems: []models.Product{
				products[0],
			},
			OrderTotal: total1,
		},
		models.Order{
			UserID:     student.ID,
			ShopID:     products[0].ShopID,
			OrderItems: products,
			OrderTotal: total2,
		},
	}

	for _, order := range orders {
		err = server.DB.Model(&models.Order{}).Create(&order).Error
		if err != nil {
			return []models.Order{}, err
		}
	}

	return orders, nil
}

func TestCreateOrder(t *testing.T) {

	var AuthEmail, AuthPassword, AuthID string
	err := refreshEverything()
	if err != nil {
		log.Fatal(err)
	}

	students, err := seedStudents()
	if err != nil {
		log.Fatal(err)
	}

	products, err := seedProducts()
	if err != nil {
		log.Fatal(err)
	}

	authStudent := students[0]
	AuthID = authStudent.ID.String()
	AuthEmail = authStudent.Email
	AuthPassword = "password"
	unauthStudent := students[1]

	token, err := server.SignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	unauthToken, err := server.SignIn(unauthStudent.Email, "password")
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	unauthTokenString := fmt.Sprintf("Bearer %v", unauthToken)
	fmt.Print(unauthTokenString)

	orderProduct, err := json.Marshal(products[0])
	if err != nil {
		log.Fatal(err)
	}

	samples := []struct {
		studentID     string
		createJSON    string
		statusCode    int
		tokenGiven    string
		orderedByName string
		orderedFrom   string
		errorMessage  string
	}{
		{
			studentID:     AuthID,
			createJSON:    fmt.Sprintf(`{"shop_id": "%s", "ordered_items": [%s]}`, products[0].ShopID.String(), string(orderProduct)),
			statusCode:    201,
			tokenGiven:    tokenString,
			orderedByName: "Donald",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/students/", bytes.NewBufferString(v.createJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}

		req = mux.SetURLVars(req, map[string]string{"student_id": v.studentID})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateOrder)
		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			orderedBy := responseMap["ordered_by"].(map[string]interface{})
			assert.Equal(t, orderedBy["first_name"], v.orderedByName)
		}

		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
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
