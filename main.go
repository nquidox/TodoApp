package main

import (
	"github.com/joho/godotenv"
	"log"
	"todoApp/config"
	"todoApp/db"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	c := config.New()
	db.ConnectToDB(c)
	server := NewApiServer(c.Config.HTTPHost, c.Config.HTTPPort)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
