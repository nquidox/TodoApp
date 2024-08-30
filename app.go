package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"todoApp/api/todoList"
	"todoApp/api/user"
	"todoApp/config"
	"todoApp/types"
)

type todoApp struct {
	dbWorker   types.DatabaseWorker
	authWorker types.AuthWorker
	salt       []byte
	server     *ApiServer
	router     *http.ServeMux
	config     *config.Config
}

func (t *todoApp) Init() error {
	user.Init(&user.Service{
		DbWorker:   t.dbWorker,
		AuthWorker: t.authWorker,
		Salt:       t.salt,
		Router:     t.router,
		Config:     t.config,
	})

	todoList.Init(&todoList.Service{
		DbWorker:   t.dbWorker,
		AuthWorker: t.authWorker,
		Router:     t.router,
	})

	return nil
}

func (t *todoApp) Run() error {
	if err := t.server.Run(t.router); err != nil {
		log.Fatal(err)
	}
	return nil
}
