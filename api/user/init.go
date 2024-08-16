package user

import (
	log "github.com/sirupsen/logrus"
	"todoApp/api/service"
	"todoApp/types"
)

var Worker types.DatabaseWorker

func Init(wrk types.DatabaseWorker) {
	var err error

	Worker = wrk

	err = wrk.InitTable(&User{})
	if err != nil {
		log.Fatal(service.TableInitErr, err)
	}

	err = wrk.InitTable(&Session{})
	if err != nil {
		log.Fatal(service.TableInitErr, err)
	}
}
