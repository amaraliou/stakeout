package handlers

import "github.com/amaraliou/apetitoso/middlewares"

func (server *Server) initializeRoutes() {

	// Look into splitting the routes into multiple files

	// /api/v1 prefix
	server.Router.PathPrefix("/api/v1") //.Subrouter()

	// Home route
	server.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(server.Home)).Methods("GET")

	// Login Route
	server.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(server.Login)).Methods("POST")
	server.Router.HandleFunc("/admins/login", middlewares.SetMiddlewareJSON(server.AdminLogin)).Methods("POST")

	// Students routes
	server.Router.HandleFunc("/students", middlewares.SetMiddlewareJSON(server.CreateStudent)).Methods("POST")
	server.Router.HandleFunc("/students", middlewares.SetMiddlewareJSON(server.GetStudents)).Methods("GET")
	server.Router.HandleFunc("/students/{id}", middlewares.SetMiddlewareJSON(server.GetStudentByID)).Methods("GET")
	server.Router.HandleFunc("/students/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.UpdateStudent))).Methods("PUT")
	server.Router.HandleFunc("/students/{id}", middlewares.SetMiddlewareAuthentication(server.DeleteStudent)).Methods("DELETE")

	// Admin routes
	server.Router.HandleFunc("/admins", middlewares.SetMiddlewareJSON(server.CreateAdmin)).Methods("POST")
	server.Router.HandleFunc("/admins", middlewares.SetMiddlewareJSON(server.GetAdmins)).Methods("GET") // Add additional Auth permissions where owners of the systems can only do this
	server.Router.HandleFunc("/admins/{id}", middlewares.SetMiddlewareJSON(server.GetAdminByID)).Methods("GET")
	server.Router.HandleFunc("/admins/{id}", middlewares.SetMiddlewareAdminAuthentication(middlewares.SetMiddlewareJSON(server.UpdateAdmin))).Methods("PUT")
	server.Router.HandleFunc("/admins/{id}", middlewares.SetMiddlewareAuthentication(server.DeleteAdmin)).Methods("DELETE")

	// Shop routes
	server.Router.HandleFunc("/admins/{admin_id}/shops", middlewares.SetMiddlewareAdminAuthentication(middlewares.SetMiddlewareJSON(server.CreateShop))).Methods("POST")
	server.Router.HandleFunc("/shops", middlewares.SetMiddlewareAuthentication(server.GetShops)).Methods("GET")
	server.Router.HandleFunc("/shops/{id}", middlewares.SetMiddlewareJSON(server.GetShopByID)).Methods("GET")
	server.Router.HandleFunc("/admins/{admin_id}/shops/{shop_id}", middlewares.SetMiddlewareAdminAuthentication(middlewares.SetMiddlewareJSON(server.UpdateShop))).Methods("PUT")
	server.Router.HandleFunc("/admins/{admin_id}/shops/{shop_id}", middlewares.SetMiddlewareAdminAuthentication(middlewares.SetMiddlewareJSON(server.DeleteShop))).Methods("DELETE")
}
