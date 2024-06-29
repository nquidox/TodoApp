package main

import (
	"github.com/joho/godotenv"
	"log"
	"todoApp/api/todoList"
	"todoApp/api/user"
	"todoApp/config"
	"todoApp/db"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	c := config.New()

	user.SALT = []byte("hglI##ERgf9D)9e5v_*ZqS=H4JN9fFAu")

	DB := db.ConnectToDB(c)
	user.Init(DB)
	todoList.Init(DB)

	server := NewApiServer(c.Config.HTTPHost, c.Config.HTTPPort)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
