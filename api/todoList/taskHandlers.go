package todoList

import (
	"github.com/google/uuid"
	"io"
	"net/http"
	"todoApp/api/service"
)

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusBadRequest,
			Messages:   "Error parsing id: " + err.Error(),
			Data:       "",
		})
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusInternalServerError,
			Messages:   "Error reading body: " + err.Error(),
			Data:       "",
		})
		return
	}

	task := Task{TodoListId: id}

	err = service.DeserializeJSON(data, &task)
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusInternalServerError,
			Messages:   "JSON read error: " + err.Error(),
			Data:       "",
		})
		return
	}

	err = task.validateTitle()
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusBadRequest,
			Messages:   "Validation error: " + err.Error(),
			Data:       "",
		})
		return
	}

	err = task.Create()
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusInternalServerError,
			Messages:   "Create task error: " + err.Error(),
			Data:       "",
		})
		return
	}

	service.ServerResponse(w, task)
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	q := r.URL.Query()
	count := validateQueryInt(q.Get("count"), 10)
	page := validateQueryInt(q.Get("page"), 1)

	id, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusBadRequest,
			Messages:   "Error parsing id: " + err.Error(),
			Data:       "",
		})
		return
	}

	tasks := Task{TodoListId: id}
	read, err := tasks.Read(count, page)
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusInternalServerError,
			Messages:   "Tasks read error: " + err.Error(),
			Data:       "",
		})
		return
	}

	service.ServerResponse(w, read)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	listId, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusBadRequest,
			Messages:   "Error parsing id: " + err.Error(),
			Data:       "",
		})
		return
	}

	taskId, err := uuid.Parse(r.PathValue("taskId"))
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusBadRequest,
			Messages:   "Error parsing id: " + err.Error(),
			Data:       "",
		})
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusBadRequest,
			Messages:   "Error reading body: " + err.Error(),
			Data:       "",
		})
		return
	}

	task := Task{TaskId: taskId, TodoListId: listId}

	err = service.DeserializeJSON(data, &task)
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusInternalServerError,
			Messages:   "JSON read error: " + err.Error(),
			Data:       "",
		})
		return
	}

	err = task.validateTitle()
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusBadRequest,
			Messages:   "Validation error: " + err.Error(),
			Data:       "",
		})
		return
	}

	err = task.Update()
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusInternalServerError,
			Messages:   "Task update error: " + err.Error(),
			Data:       "",
		})
		return
	}

	service.ServerResponse(w, task)
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	listId, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusBadRequest,
			Messages:   "Error parsing id: " + err.Error(),
			Data:       "",
		})
		return
	}

	taskId, err := uuid.Parse(r.PathValue("taskId"))
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusBadRequest,
			Messages:   "Error parsing id: " + err.Error(),
			Data:       "",
		})
		return
	}

	t := Task{TodoListId: listId, TaskId: taskId}
	err = t.Delete()
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusInternalServerError,
			Messages:   "Delete error: " + err.Error(),
			Data:       "",
		})
		return
	}

	service.ServerResponse(w, service.ErrorResponse{
		ResultCode: 0,
		ErrorCode:  http.StatusOK,
		Messages:   "",
		Data:       "",
	})
}
