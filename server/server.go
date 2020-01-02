package server

import (
	"fmt"
	"log"
	"os"

	"github.com/amaraliou/apetitoso/handlers"
	"github.com/amaraliou/apetitoso/utils"
	"github.com/joho/godotenv"
)

var server = handlers.Server{}

// Run -> start server
func Run() {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	server.Initialize(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	utils.Load(server.DB)

	server.Run(":8080")
}