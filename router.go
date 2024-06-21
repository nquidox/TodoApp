package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"todoApp/db"
)

type ApiServer struct {
	Addr string
}

type Route struct {
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

type Router struct {
	routes []Route
}

func NewApiServer(host, port string) *ApiServer {
	return &ApiServer{Addr: host + ":" + port}
}

func NewRouter() *Router {
	return &Router{}
}

func (s *ApiServer) Run() error {
	router := NewRouter()
	prefix := "/api/v1"

	server := &http.Server{Addr: s.Addr, Handler: router}

	router.HandleFunc("GET", fmt.Sprintf("%s", prefix), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "API v1 is ready.")
	})

	router.HandleFunc("GET", fmt.Sprintf("%s/createuser", prefix), func(w http.ResponseWriter, r *http.Request) {
		err := db.CreateUser([]byte(`{"id":1,"nickname":"vasya_pupkin","name":"Vasya","surname":"Pupkin","email":"vasya.pupkin@example.com","password":"securepassword123"}`))

		if err != nil {
			fmt.Fprint(w, "Error creating user")
		} else {
			fmt.Fprint(w, "User created")
		}

	})

	log.Printf("Starting server on %s", s.Addr)
	return server.ListenAndServe()
}

func (r *Router) HandleFunc(method, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	r.routes = append(r.routes, Route{
		Method:  method,
		Pattern: pattern,
		Handler: http.HandlerFunc(handler),
	})
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		if req.Method == route.Method && strings.HasPrefix(req.URL.Path, route.Pattern) {
			if req.URL.Path == route.Pattern || req.URL.Path == route.Pattern+"/" {
				route.Handler.ServeHTTP(w, req)
				return
			}
		}
	}
	http.NotFound(w, req)
}
