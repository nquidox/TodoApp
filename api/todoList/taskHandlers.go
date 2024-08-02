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
		service.BadRequestResponse(w, service.ParseErr, err)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		service.BadRequestResponse(w, service.BodyReadErr, err)
		return
	}

	task := Task{TodoListId: id}

	err = service.DeserializeJSON(data, &task)
	if err != nil {
		service.UnprocessableEntity(w, service.JSONDeserializingErr, err)
		return
	}

	err = task.validateTitle()
	if err != nil {
		service.BadRequestResponse(w, service.ValidationErr, err)
		return
	}

	err = task.Create()
	if err != nil {
		service.InternalServerErrorResponse(w, service.TaskCreateErr, err)
		return
	}

	service.OkResponse(w, task)
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	q := r.URL.Query()
	count := validateQueryInt(q.Get("count"), 10)
	page := validateQueryInt(q.Get("page"), 1)

	id, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		service.BadRequestResponse(w, service.ParseErr, err)
		return
	}

	tasks := Task{TodoListId: id}
	read, err := tasks.Read(count, page)
	if err != nil {
		service.InternalServerErrorResponse(w, service.TaskReadErr, err)
		return
	}

	service.OkResponse(w, read)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	listId, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		service.BadRequestResponse(w, service.ParseErr, err)
		return
	}

	taskId, err := uuid.Parse(r.PathValue("taskId"))
	if err != nil {
		service.BadRequestResponse(w, service.ParseErr, err)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		service.BadRequestResponse(w, service.BodyReadErr, err)
		return
	}

	task := Task{TaskId: taskId, TodoListId: listId}

	err = service.DeserializeJSON(data, &task)
	if err != nil {
		service.UnprocessableEntity(w, service.JSONDeserializingErr, err)
		return
	}

	err = task.validateTitle()
	if err != nil {
		service.BadRequestResponse(w, service.ValidationErr, err)
		return
	}

	err = task.Update()
	if err != nil {
		service.InternalServerErrorResponse(w, service.TaskUpdateErr, err)
		return
	}

	service.OkResponse(w, task)
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	listId, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		service.BadRequestResponse(w, service.ParseErr, err)
		return
	}

	taskId, err := uuid.Parse(r.PathValue("taskId"))
	if err != nil {
		service.BadRequestResponse(w, service.ParseErr, err)
		return
	}

	t := Task{TodoListId: listId, TaskId: taskId}
	err = t.Delete()
	if err != nil {
		service.InternalServerErrorResponse(w, service.TaskDeleteErr, err)
		return
	}

	service.OkResponse(w, service.DefaultResponse{
		ResultCode: 0,
		HttpCode:   http.StatusOK,
		Messages:   "",
		Data:       "",
	})
}
