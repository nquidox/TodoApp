package todoList

import (
	"log"
	"todoApp/api/service"
	"todoApp/types"
)

type dbWorker types.DatabaseWorker
type authWorker types.AuthWorker
type authUser struct{ types.AuthUser }

var (
	dbw types.DatabaseWorker
	aw  types.AuthWorker
)

func Init(dbWrk dbWorker, authWkr authWorker) {
	var err error
	dbw = dbWrk
	aw = authWkr

	err = dbWrk.InitTable(&TodoList{})
	if err != nil {
		log.Fatal(service.TableInitErr, err)
	}

	err = dbWrk.InitTable(&Task{})
	if err != nil {
		log.Fatal(service.TableInitErr, err)
	}
}
