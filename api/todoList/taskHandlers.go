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
		serverResponse(w, ErrorResponse{
			ResultCode: http.StatusBadRequest,
			Messages:   "ID error: " + err.Error(),
			Data:       "",
		})
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		serverResponse(w, ErrorResponse{
			ResultCode: http.StatusInternalServerError,
			Messages:   "Internal server error: " + err.Error(),
			Data:       "",
		})
		return
	}

	title := RequestTitle{}

	err = service.DeserializeJSON(data, title)
	if err != nil {
		serverResponse(w, ErrorResponse{
			ResultCode: http.StatusInternalServerError,
			Messages:   "Internal server error: " + err.Error(),
			Data:       "",
		})
		return
	}

	task := Task{}

	err = task.Create(id, title.Title)
	if err != nil {
		serverResponse(w, ErrorResponse{
			ResultCode: http.StatusInternalServerError,
			Messages:   "Internal server error: " + err.Error(),
			Data:       "",
		})
		return
	}

	serverResponse(w, task)
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	q := r.URL.Query()
	count := validateQueryInt(q.Get("count"), 10)
	page := validateQueryInt(q.Get("page"), 1)

	id, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		serverResponse(w, ErrorResponse{
			ResultCode: http.StatusBadRequest,
			Messages:   "List ID error: " + err.Error(),
			Data:       "",
		})
		return
	}

	tasks := Task{}
	read, err := tasks.Read(id, count, page)
	if err != nil {
		serverResponse(w, ErrorResponse{
			ResultCode: http.StatusInternalServerError,
			Messages:   "Read error: " + err.Error(),
			Data:       "",
		})
		return
	}

	serverResponse(w, read)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	listId, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		serverResponse(w, ErrorResponse{
			ResultCode: http.StatusBadRequest,
			Messages:   "List ID error: " + err.Error(),
			Data:       "",
		})
		return
	}

	taskId, err := uuid.Parse(r.PathValue("taskId"))
	if err != nil {
		serverResponse(w, ErrorResponse{
			ResultCode: http.StatusBadRequest,
			Messages:   "Task ID error: " + err.Error(),
			Data:       "",
		})
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		serverResponse(w, ErrorResponse{
			ResultCode: http.StatusBadRequest,
			Messages:   "Body parse error: " + err.Error(),
			Data:       "",
		})
		return
	}

	task := Task{}

	err = service.DeserializeJSON(data, &task)
	if err != nil {
		serverResponse(w, ErrorResponse{
			ResultCode: http.StatusInternalServerError,
			Messages:   "Task parse error: " + err.Error(),
			Data:       "",
		})
		return
	}

	result, err := task.Update(listId, taskId)
	if err != nil {
		serverResponse(w, ErrorResponse{
			ResultCode: http.StatusInternalServerError,
			Messages:   "Task update error: " + err.Error(),
			Data:       "",
		})
		return
	}

	serverResponse(w, result)
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: refactor error write
	listId, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		serverResponse(w, ErrorResponse{
			ResultCode: http.StatusBadRequest,
			Messages:   "List ID error: " + err.Error(),
			Data:       "",
		})
		return
	}

	taskId, err := uuid.Parse(r.PathValue("taskId"))
	if err != nil {
		serverResponse(w, ErrorResponse{
			ResultCode: http.StatusBadRequest,
			Messages:   "Task ID error: " + err.Error(),
			Data:       "",
		})
		return
	}

	t := Task{}
	err = t.Delete(listId, taskId)
	if err != nil {
		serverResponse(w, ErrorResponse{
			ResultCode: http.StatusInternalServerError,
			Messages:   "Delete error: " + err.Error(),
			Data:       "",
		})
		return
	}
}
