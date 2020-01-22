package handlers

import "net/http"

// CreateProduct -> handles POST /api/v1/shops/<shop_id:uuid>/products/
// Might need admin id
func (server *Server) CreateProduct(writer http.ResponseWriter, request *http.Request) {

}

// GetProducts -> handles GET /api/v1/products/
func (server *Server) GetProducts(writer http.ResponseWriter, request *http.Request) {

}

// GetProductByID -> handles Get /api/v1/products/<id:uuid>
func (server *Server) GetProductByID(writer http.ResponseWriter, request *http.Request) {

}

// GetProductsByShop -> handles GET /api/v1/shops/<shop_id:uuid>/products/
func (server *Server) GetProductsByShop(writer http.ResponseWriter, request *http.Request) {

}

// UpdateProduct -> handles PUT /api/v1/shops/<shop_id:uuid>/products/<id:uuid>
// Might need admin id
func (server *Server) UpdateProduct(writer http.ResponseWriter, request *http.Request) {

}

// DeleteProduct -> handles DELETE /api/v1/shops/<shop_id:uuid>/products/<id:uuid>
// Might need admin id
func (server *Server) DeleteProduct(writer http.ResponseWriter, request *http.Request) {

}
