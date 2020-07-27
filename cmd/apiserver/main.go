package main

import (
	"log"

	"github.com/http-rest-api/internal/app/apiserver"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("configs/apiserver.env"); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	config := apiserver.NewConfig()

	s := apiserver.New(config)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
