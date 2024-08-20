package main

import (
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"path"
	"runtime"
	"todoApp/api/todoList"
	"todoApp/api/user"
	"todoApp/config"
	"todoApp/db"
	"todoApp/types"
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
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf(" %s:%d", filename, f.Line)
		},
	})

	user.SALT = []byte("hglI##ERgf9D)9e5v_*ZqS=H4JN9fFAu")

	var dbWorker types.DatabaseWorker = &db.DB{Connection: db.Connect(c)}
	var authWorker types.AuthWorker = &user.AuthService{}

	user.Init(dbWorker)
	todoList.Init(dbWorker, authWorker)

	server := NewApiServer(c.Config.HTTPHost, c.Config.HTTPPort)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
