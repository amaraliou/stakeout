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
)

// CreateShop -> handles POST /api/v1/admins/<admin_id:uuid>/shops/
func (server *Server) CreateShop(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	adminID := vars["admin_id"]

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

	isAdmin, err := auth.IsAdminToken(request)
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	if !isAdmin {
		responses.ERROR(writer, http.StatusUnauthorized, errors.New("Unauthorized: This is not an admin token"))
		return
	}

	tokenID, err := auth.ExtractTokenAdminID(request)
	if err != nil {
		responses.ERROR(writer, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if tokenID != adminID {
		responses.ERROR(writer, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
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

// GetShops -> handles GET /api/v1/shops/
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

// UpdateShop -> handles PUT /api/v1/admins/<admin_id:uuid>/shop/<shop_id:uuid>
func (server *Server) UpdateShop(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	shopID := vars["shop_id"]
	adminID := vars["admin_id"]
	shop := models.Shop{}
	admin := models.Admin{}

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

	isAdmin, err := auth.IsAdminToken(request)
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	if !isAdmin {
		responses.ERROR(writer, http.StatusUnauthorized, errors.New("Unauthorized: This is not an admin token"))
		return
	}

	tokenID, err := auth.ExtractTokenAdminID(request)
	if err != nil {
		responses.ERROR(writer, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if tokenID != adminID {
		responses.ERROR(writer, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	currentAdmin, err := admin.FindAdminByID(server.DB, adminID)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	if currentAdmin.ShopID.String() != shopID {
		responses.ERROR(writer, http.StatusUnauthorized, errors.New("Unauthorized: You are not the admin for this shop"))
		return
	}

	updatedShop, err := shop.UpdateShop(server.DB, shopID)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(writer, http.StatusOK, updatedShop)
}

// DeleteShop -> handles DELETE /api/v1/admins/<admin_id:uuid>/shops/<shop_id:uuid>
func (server *Server) DeleteShop(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	shopID := vars["shop_id"]
	adminID := vars["admin_id"]
	shop := models.Shop{}
	admin := models.Admin{}

	isAdmin, err := auth.IsAdminToken(request)
	if err != nil {
		responses.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	if !isAdmin {
		responses.ERROR(writer, http.StatusUnauthorized, errors.New("Unauthorized: This is not an admin token"))
		return
	}

	tokenID, err := auth.ExtractTokenAdminID(request)
	if err != nil {
		responses.ERROR(writer, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if tokenID != adminID {
		responses.ERROR(writer, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	currentAdmin, err := admin.FindAdminByID(server.DB, adminID)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	if currentAdmin.ShopID.String() != shopID {
		responses.ERROR(writer, http.StatusUnauthorized, errors.New("Unauthorized: You are not the admin for this shop"))
		return
	}

	_, err = shop.DeleteShop(server.DB, shopID)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	writer.Header().Set("Entity", fmt.Sprintf("%s", shopID))
	responses.JSON(writer, http.StatusNoContent, "")
}
