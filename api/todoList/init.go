package todoList

import (
	"log"
	"net/http"
	"todoApp/api/service"
	"todoApp/types"
)

type (
	dbWorker types.DatabaseWorker
	authUser struct{ types.AuthUser }
)

type Service struct {
	DbWorker   types.DatabaseWorker
	AuthWorker types.AuthWorker
	Router     *http.ServeMux
}

func Init(s *Service) {
	var err error

	err = s.DbWorker.InitTable(&TodoList{})
	if err != nil {
		log.Fatal(service.TableInitErr, err)
	}

	err = s.DbWorker.InitTable(&Task{})
	if err != nil {
		log.Fatal(service.TableInitErr, err)
	}

	addRoutes(s)
}
