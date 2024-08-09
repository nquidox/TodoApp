package user

import (
	log "github.com/sirupsen/logrus"
	"todoApp/api/service"
)

type DatabaseWorker interface {
	InitTable(model any) error
	CreateRecord(model any) error
	ReadOneRecord(model any, params map[string]any) error
	ReadManyRecords(model any) error
	UpdateRecord(model any, params map[string]any) error
	DeleteRecord(model any, params map[string]any) error
}

var Worker DatabaseWorker

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
