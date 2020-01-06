package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/amaraliou/apetitoso/models"
	"github.com/amaraliou/apetitoso/responses"
	"github.com/gorilla/mux"
)

// CreateShop -> handles POST /api/v1/shop/
func (server *Server) CreateShop(writer http.ResponseWriter, request *http.Request) {

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
	}

	shop := models.Shop{}
	err = json.Unmarshal(body, &shop)
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	err = shop.Validate("create")
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	shopCreated, err := shop.CreateShop(server.DB)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	writer.Header().Set("Location", fmt.Sprintf("%s%s/%s", request.Host, request.RequestURI, shopCreated.ID.String()))
	responses.JSON(writer, http.StatusCreated, shopCreated)
}

// GetShops -> handles GET /api/v1/shop/
func (server *Server) GetShops(writer http.ResponseWriter, request *http.Request) {

	shop := models.Shop{}
	shops, err := shop.FindAllShops(server.DB)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(writer, http.StatusOK, shops)
}

// GetShopByID -> handles GET /api/v1/shop/<id:uuid>
func (server *Server) GetShopByID(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	shop := models.Shop{}
	shopRetrieved, err := shop.FindShopByID(server.DB, vars["id"])
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(writer, http.StatusOK, shopRetrieved)
}

// UpdateShop -> handles PUT /api/v1/shop/<id:uuid>
func (server *Server) UpdateShop(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	shopID := vars["id"]
	shop := models.Shop{}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
	}

	err = json.Unmarshal(body, &shop)
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	err = shop.Validate("")
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	// Verify who the current user/admin is to check permissions (maybe change the endpoint to /api/v1/admin/<admin_id:uuid>/shop/<shop_id:uuid>)

	updatedShop, err := shop.UpdateShop(server.DB, shopID)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(writer, http.StatusOK, updatedShop)
}

// DeleteShop -> handles DELETE /api/v1/shop/<id:uuid> (make sure products are deleted as well)
func (server *Server) DeleteShop(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	shopID := vars["id"]
	shop := models.Shop{}

	// Verify who the current user/admin is to check permissions (maybe change the endpoint to /api/v1/admin/<admin_id:uuid>/shop/<shop_id:uuid>)

	_, err := shop.DeleteShop(server.DB, shopID)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	writer.Header().Set("Entity", fmt.Sprintf("%s", shopID))
	responses.JSON(writer, http.StatusNoContent, "")
}
