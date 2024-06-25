package main

import (
	"fmt"
	"log"
	"net/http"
	"todoApp/api/todoList"
	"todoApp/api/user"
)

type ApiServer struct {
	Addr string
}

func NewApiServer(host, port string) *ApiServer {
	return &ApiServer{Addr: host + ":" + port}
}

func (s *ApiServer) Run() error {
	router := http.NewServeMux()
	server := &http.Server{Addr: s.Addr, Handler: router}

	router.HandleFunc("GET /api/v1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "API v1 is ready.")
	})

	user.AddRoutes(router)
	todoList.AddRoutes(router)

	log.Printf("Starting server on %s", s.Addr)
	return server.ListenAndServe()
}
