package main

import (
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"todoApp/api/todoList"
	"todoApp/api/user"
	_ "todoApp/docs"
)

type ApiServer struct {
	Addr string
}

func NewApiServer(host, port string) *ApiServer {
	return &ApiServer{Addr: host + ":" + port}
}

func (s *ApiServer) Run() error {
	router := http.NewServeMux()
	server := &http.Server{Addr: s.Addr, Handler: corsMiddleware(router)}

	router.HandleFunc("GET /api/v1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "API v1 is ready.")
	})

	router.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	user.AddRoutes(router)
	todoList.AddRoutes(router)

	log.Printf("Starting server on %s", s.Addr)
	return server.ListenAndServe()
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
