package handlers

import "github.com/amaraliou/apetitoso/middlewares"

func (server *Server) initializeRoutes() {

	// /api/v1 prefix
	server.Router.PathPrefix("/api/v1")

	// Home route
	server.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(server.Home)).Methods("GET")

	// Login Route
	server.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(server.Login)).Methods("POST")

	// Students routes
	server.Router.HandleFunc("/students", middlewares.SetMiddlewareJSON(server.CreateStudent)).Methods("POST")
	server.Router.HandleFunc("/students", middlewares.SetMiddlewareJSON(server.GetStudents)).Methods("GET")
	server.Router.HandleFunc("/students/{id}", middlewares.SetMiddlewareJSON(server.GetStudentByID)).Methods("GET")
	server.Router.HandleFunc("/students/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.UpdateStudent))).Methods("PUT")
	server.Router.HandleFunc("/students/{id}", middlewares.SetMiddlewareAuthentication(server.DeleteStudent)).Methods("DELETE")

	// Admin routes
}
