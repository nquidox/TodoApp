package user

import (
	"net/http"
)

func AddRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /api/v1/user", CreateUserHandler)
	router.HandleFunc("GET /api/v1/user/{id}", ReadUserHandler)
	router.HandleFunc("PUT /api/v1/user/{id}", UpdateUserHandler)
	router.HandleFunc("DELETE /api/v1/user/{id}", DeleteUserHandler)

	router.HandleFunc("POST /api/v1/login", LoginHandler)
	router.HandleFunc("GET /api/v1/logout", LogoutHandler)

	router.HandleFunc("GET /api/v1/me", MeHandler)
}
