package user

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"todoApp/api/service"
	"todoApp/config"
	"todoApp/types"
)

type dbWorker types.DatabaseWorker

type Service struct {
	DbWorker   types.DatabaseWorker
	AuthWorker types.AuthWorker
	Salt       []byte
	Router     *http.ServeMux
	Config     *config.Config
}

func Init(s *Service) {
	var err error

	err = s.DbWorker.InitTable(&User{})
	if err != nil {
		log.Fatal(service.TableInitErr, err)
	}

	err = s.DbWorker.InitTable(&Session{})
	if err != nil {
		log.Fatal(service.TableInitErr, err)
	}

	addRoutes(s)
}
