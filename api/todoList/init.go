package todoList

import (
	"log"
	"todoApp/api/service"
	"todoApp/types"
)

var Worker types.DatabaseWorker

func Init(wrk types.DatabaseWorker) {
	var err error
	Worker = wrk

	err = wrk.InitTable(&TodoList{})
	if err != nil {
		log.Fatal(service.TableInitErr, err)
	}

	err = wrk.InitTable(&Task{})
	if err != nil {
		log.Fatal(service.TableInitErr, err)
	}
}
