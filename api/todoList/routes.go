package todoList

import (
	"net/http"
)

func AddRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /api/v1/todo-lists", CreateListHandler)
	router.HandleFunc("GET /api/v1/todo-lists", GetAllListsHandler)
	router.HandleFunc("PUT /api/v1/todo-lists/{id}", UpdateListHandler)
	router.HandleFunc("DELETE /api/v1/todo-lists/{id}", DeleteListHandler)
}
