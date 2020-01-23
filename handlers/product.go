package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/amaraliou/apetitoso/auth"
	"github.com/amaraliou/apetitoso/models"
	"github.com/amaraliou/apetitoso/responses"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// CreateProduct -> handles POST /api/v1/shops/<shop_id:uuid>/products/
// Might need admin id
func (server *Server) CreateProduct(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	shopID := vars["shop_id"]

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
	}

	product := models.Product{}
	err = json.Unmarshal(body, &product)
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	isAdmin, err := auth.IsAdminToken(request)
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	if !isAdmin {
		responses.ERROR(writer, http.StatusUnauthorized, errors.New("Unauthorized: This is not an admin token"))
		return
	}

	err = product.Validate("create")
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	shopUUID, err := uuid.FromString(shopID)
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	product.ShopID = shopUUID

	productCreated, err := product.CreateProduct(server.DB)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	writer.Header().Set("Location", fmt.Sprintf("%s%s/%s", request.Host, request.RequestURI, productCreated.ID.String()))
	responses.JSON(writer, http.StatusCreated, productCreated)
}

// GetProducts -> handles GET /api/v1/products/
func (server *Server) GetProducts(writer http.ResponseWriter, request *http.Request) {

	product := models.Product{}
	products, err := product.FindAllProducts(server.DB)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(writer, http.StatusOK, products)
}

// GetProductByID -> handles Get /api/v1/products/<id:uuid>
func (server *Server) GetProductByID(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	product := models.Product{}
	productRetrieved, err := product.FindProductByID(server.DB, vars["id"])
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(writer, http.StatusOK, productRetrieved)
}

// GetProductsByShop -> handles GET /api/v1/shops/<shop_id:uuid>/products/
func (server *Server) GetProductsByShop(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	product := models.Product{}
	products, err := product.FindAllProductsByShop(server.DB, vars["shop_id"])
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(writer, http.StatusOK, products)
}

// UpdateProduct -> handles PUT /api/v1/shops/<shop_id:uuid>/products/<id:uuid>
// Might need admin id
func (server *Server) UpdateProduct(writer http.ResponseWriter, request *http.Request) {

}

// DeleteProduct -> handles DELETE /api/v1/shops/<shop_id:uuid>/products/<id:uuid>
// Might need admin id
func (server *Server) DeleteProduct(writer http.ResponseWriter, request *http.Request) {

}