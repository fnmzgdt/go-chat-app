package main

import (
	"log"
	"project/pkg/router"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
  	if err != nil {
    	log.Printf("Error loading .env file")
  	}
	router.StartServer()
}
