package user

import (
	log "github.com/sirupsen/logrus"
	"todoApp/api/service"
	"todoApp/types"
)

var Worker types.DatabaseWorker

func Init() {
	var err error

	err = Worker.InitTable(&User{})
	if err != nil {
		log.Fatal(service.TableInitErr, err)
	}

	err = Worker.InitTable(&Session{})
	if err != nil {
		log.Fatal(service.TableInitErr, err)
	}
}
