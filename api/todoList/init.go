package todoList

import (
	"log"
	"todoApp/api/service"
	"todoApp/types"
)

var (
	dbWorker   types.DatabaseWorker
	authWorker types.AuthWorker
)

func Init(dbWrk types.DatabaseWorker, authWkr types.AuthWorker) {
	var err error
	dbWorker = dbWrk
	authWorker = authWkr

	err = dbWrk.InitTable(&TodoList{})
	if err != nil {
		log.Fatal(service.TableInitErr, err)
	}

	err = dbWrk.InitTable(&Task{})
	if err != nil {
		log.Fatal(service.TableInitErr, err)
	}
}
