package main

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"todoApp/api/todoList"
	"todoApp/api/user"
	"todoApp/config"
	"todoApp/db"
)

// @title						TODO App API
// @version					1.0
// @description				This is an educational TODO App API, clone of IT-Incubator's API written in Go.
//
// @contact.name				API Support
// @contact.url				https://t.me/rekasawak
//
// @license.name				MIT
// @license.url				https://mit-license.org/
//
// @host						localhost:9000
// @BasePath					/api/v1
//
// @securityDefinitions.basic	BasicAuth
// @in							header
// @name						token
//
// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	c := config.New()

	log.SetLevel(appSetLogLevel(c.Config.AppLogLevel))
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "02-01-2006 15:04:05",
	})

	user.SALT = []byte("hglI##ERgf9D)9e5v_*ZqS=H4JN9fFAu")

	DB := db.ConnectToDB(c)
	user.Init(DB)
	todoList.Init(DB)

	server := NewApiServer(c.Config.HTTPHost, c.Config.HTTPPort)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
