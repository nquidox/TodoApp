package main

import (
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"net/http"
	"path"
	"runtime"
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
	var err error
	if err = godotenv.Load(); err != nil {
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

	app := todoApp{
		dbWorker:   &db.DB{Connection: db.Connect(c)},
		authWorker: &user.AuthService{},
		salt:       []byte("hglI##ERgf9D)9e5v_*ZqS=H4JN9fFAu"),
		server:     NewApiServer(c.Config.HTTPHost, c.Config.HTTPPort),
		router:     http.NewServeMux(),
		config:     c,
		cors:       config.CorsConfig(),
	}

	err = app.Init()
	if err != nil {
		log.WithError(err).Fatal("Error initializing server")
	}

	err = app.Run()
	if err != nil {
		log.WithError(err).Fatal("Error starting server")
	}
}
