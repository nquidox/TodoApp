package user

import (
	log "github.com/sirupsen/logrus"
	"todoApp/api/service"
	"todoApp/types"
)

type dbWorker types.DatabaseWorker

var dbw dbWorker

func Init(wrk dbWorker) {
	var err error

	dbw = wrk

	err = wrk.InitTable(&User{})
	if err != nil {
		log.Fatal(service.TableInitErr, err)
	}

	err = wrk.InitTable(&Session{})
	if err != nil {
		log.Fatal(service.TableInitErr, err)
	}
}
