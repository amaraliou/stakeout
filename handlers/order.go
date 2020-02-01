package handlers

import "net/http"

// CreateOrder -> handles POST /api/v1/students/<student_id:uuid>/orders/
func (server *Server) CreateOrder(writer http.ResponseWriter, request *http.Request) {

}

// GetAllOrders -> handles GET /api/v1/orders/
func (server *Server) GetAllOrders(writer http.ResponseWriter, request *http.Request) {

}

// GetAllOrdersByStudent -> handles GET /api/v1/students/<student_id:uuid>/orders/
func (server *Server) GetAllOrdersByStudent(writer http.ResponseWriter, request *http.Request) {

}

// GetAllOrdersByShop -> handles GET /api/v1/shops/<shop_id:uuid>/orders/
func (server *Server) GetAllOrdersByShop(writer http.ResponseWriter, request *http.Request) {

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
