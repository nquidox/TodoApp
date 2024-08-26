package todoList

func addRoutes(s *Service) {

	createListHandler := createListFunc(s)
	s.Router.HandleFunc("POST /api/v1/todo-lists", createListHandler)

	getAllListsHandler := getAllListsFunc(s)
	s.Router.HandleFunc("GET /api/v1/todo-lists", getAllListsHandler)

	updateListHandler := updateListFunc(s)
	s.Router.HandleFunc("PUT /api/v1/todo-lists/{listId}", updateListHandler)

	deleteListHandler := deleteListFunc(s)
	s.Router.HandleFunc("DELETE /api/v1/todo-lists/{listId}", deleteListHandler)

	createTaskHandler := createTaskFunc(s)
	s.Router.HandleFunc("POST /api/v1/todo-lists/{listId}/tasks", createTaskHandler)

	getTaskHandler := getTaskFunc(s)
	s.Router.HandleFunc("GET /api/v1/todo-lists/{listId}/tasks", getTaskHandler)

	updateTaskHandler := updateTaskFunc(s)
	s.Router.HandleFunc("PUT /api/v1/todo-lists/{listId}/tasks/{taskId}", updateTaskHandler)

	deleteTaskHandler := deleteTaskFunc(s)
	s.Router.HandleFunc("DELETE /api/v1/todo-lists/{listId}/tasks/{taskId}", deleteTaskHandler)
}
