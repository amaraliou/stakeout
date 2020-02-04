package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/amaraliou/stakeout/auth"
	"github.com/amaraliou/stakeout/models"
	"github.com/amaraliou/stakeout/responses"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// CreateOrder -> handles POST /api/v1/students/<student_id:uuid>/orders/
func (server *Server) CreateOrder(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	studentID := vars["student_id"]
	student := models.Student{}
	shop := models.Shop{}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
	}

	order := models.Order{}
	err = json.Unmarshal(body, &order)
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	tokenID, err := auth.ExtractTokenID(request)
	if err != nil {
		responses.ERROR(writer, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if tokenID != studentID {
		responses.ERROR(writer, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	_, err = student.FindStudentByID(server.DB, studentID)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	studentUUID, err := uuid.FromString(studentID)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, errors.New("Invalid student UUID format"))
		return
	}

	order.UserID = studentUUID

	err = order.Validate("create")
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	_, err = shop.FindShopByID(server.DB, order.ShopID.String())
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	orderCreated, err := order.CreateOrder(server.DB)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	writer.Header().Set("Location", fmt.Sprintf("%s%s/%s", request.Host, request.RequestURI, orderCreated.ID.String()))
	responses.JSON(writer, http.StatusCreated, orderCreated)
}

// GetAllOrders -> handles GET /api/v1/orders/
func (server *Server) GetAllOrders(writer http.ResponseWriter, request *http.Request) {

	order := models.Order{}
	orders, err := order.FindAllOrders(server.DB)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(writer, http.StatusOK, orders)
}

// GetAllOrdersByStudent -> handles GET /api/v1/students/<student_id:uuid>/orders/
func (server *Server) GetAllOrdersByStudent(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	studentID := vars["student_id"]
	student := models.Student{}
	order := models.Order{}

	tokenID, err := auth.ExtractTokenID(request)
	if err != nil {
		responses.ERROR(writer, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if tokenID != studentID {
		responses.ERROR(writer, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	_, err = student.FindStudentByID(server.DB, studentID)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	orders, err := order.FindAllOrdersByStudent(server.DB, studentID)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(writer, http.StatusOK, orders)
}

// GetAllOrdersByShop -> handles GET /api/v1/shops/<shop_id:uuid>/orders/
func (server *Server) GetAllOrdersByShop(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	shopID := vars["shop_id"]
	shop := models.Shop{}
	order := models.Order{}

	// Verify that it's admin, student or owner (soon to be introduced)
	err := auth.TokenValid(request)
	if err != nil {
		responses.ERROR(writer, http.StatusUnauthorized, errors.New("Unauthorized: Admin, Student or Owner not authenticated"))
		return
	}

	_, err = shop.FindShopByID(server.DB, shopID)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	orders, err := order.FindAllOrdersByShop(server.DB, shopID)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(writer, http.StatusOK, orders)
}

// GetOrderByID -> handles GET /api/v1/orders/<id:uuid>
func (server *Server) GetOrderByID(writer http.ResponseWriter, request *http.Request) {

}

// UpdateOrder -> handles PUT /api/v1/shops/<shop_id:uuid>/orders/<order_id:uuid>
func (server *Server) UpdateOrder(writer http.ResponseWriter, request *http.Request) {

}

// DeleteOrder -> handles DELETE /api/v1/shops/<shop_id:uuid>/orders/<order_id:uuid>
func (server *Server) DeleteOrder(writer http.ResponseWriter, request *http.Request) {

}
