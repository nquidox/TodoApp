package todoList

import (
	"log"
	"todoApp/api/service"
)

type DatabaseWorker interface {
	InitTable(model any) error
	CreateRecord(model any) error
	ReadOneRecord(model any, params map[string]any) error
	ReadManyRecords(model any) error
	ReadWithPagination(model any, params map[string]any) error
	UpdateRecord(model any, params map[string]any) error
	DeleteRecord(model any, params map[string]any) error
}

var Worker DatabaseWorker

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
