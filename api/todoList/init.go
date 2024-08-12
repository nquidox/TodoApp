package todoList

import (
	"log"
	"todoApp/api/service"
	"todoApp/types"
)

var Worker types.DatabaseWorker

func Init() {
	var err error

	err = Worker.InitTable(&TodoList{})
	if err != nil {
		log.Fatal(service.TableInitErr, err)
	}

	err = Worker.InitTable(&Task{})
	if err != nil {
		log.Fatal(service.TableInitErr, err)
	}
}
