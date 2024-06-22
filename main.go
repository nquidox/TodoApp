package main

import (
	"github.com/joho/godotenv"
	"log"
	"todoApp/config"
	"todoApp/db"
	"todoApp/user"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	c := config.New()

	DB := db.ConnectToDB(c)

	user.DB = DB
	user.InitUsers()

	server := NewApiServer(c.Config.HTTPHost, c.Config.HTTPPort)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
