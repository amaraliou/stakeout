package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// Server ...
type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

// Initialize -> Function to initialize a server with Postgres given the credentials
func (server *Server) Initialize(User, Password, Port, Host, Name string) {

	var err error

	DBURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", Host, Port, User, Name, Password)
	server.DB, err = gorm.Open("postgres", DBURI)
	if err != nil {
		fmt.Printf("Cannot connect to Postgres database")
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the Postgres database")
	}

	server.DB.Debug().AutoMigrate()
	server.Router = mux.NewRouter()
	server.initializeRoutes()
}

// Run ... making my linter happy
func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
