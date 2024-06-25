package user

import (
	"fmt"
	"net/http"
)

func AddRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /api/v1/create", CreateUserHandler)
	router.HandleFunc(fmt.Sprintf("GET /api/v1/read/{id}"), ReadUserHandler)
	router.HandleFunc("UPDATE /api/v1/update/{id}", UpdateUserHandler)
	router.HandleFunc("DELETE /api/v1/delete/{id}", DeleteUserHandler)
	router.HandleFunc("POST /api/v1/login", LoginHandler)
	router.HandleFunc("GET /api/v1/logout", LogoutHandler)
}
