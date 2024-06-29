package todoList

import (
	"net/http"
)

func AddRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /api/v1/todo-lists", CreateListHandler)
	router.HandleFunc("GET /api/v1/todo-lists", GetAllListsHandler)
	router.HandleFunc("PUT /api/v1/todo-lists/{listId}", UpdateListHandler)
	router.HandleFunc("DELETE /api/v1/todo-lists/{listId}", DeleteListHandler)

	router.HandleFunc("POST /api/v1/todo-lists/{listId}/tasks", CreateTaskHandler)
	router.HandleFunc("GET /api/v1/todo-lists/{listId}/tasks", GetTaskHandler)
	router.HandleFunc("PUT /api/v1/todo-lists/{listId}/tasks/{taskId}", UpdateTaskHandler)
	router.HandleFunc("DELETE /api/v1/todo-lists/{listId}/tasks/{taskId}", DeleteTaskHandler)
}
