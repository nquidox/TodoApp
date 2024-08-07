package todoList

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"todoApp/api/service"
)

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		log.Error(service.ParseErr, err)
		service.BadRequestResponse(w, service.ParseErr, err)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(service.BodyReadErr, err)
		service.BadRequestResponse(w, service.BodyReadErr, err)
		return
	}

	task := Task{TodoListId: id}

	err = service.DeserializeJSON(data, &task)
	if err != nil {
		log.Error(service.JSONDeserializingErr, err)
		service.UnprocessableEntity(w, service.JSONDeserializingErr, err)
		return
	}

	err = task.validateTitle()
	if err != nil {
		log.Error(service.ValidationErr, err)
		service.BadRequestResponse(w, service.ValidationErr, err)
		return
	}

	err = task.Create()
	if err != nil {
		log.Error(service.TaskCreateErr, err)
		service.InternalServerErrorResponse(w, service.TaskCreateErr, err)
		return
	}

	service.OkResponse(w, task)

	log.WithFields(log.Fields{
		"id": task.ID,
	}).Info(service.TaskCreateSuccess)
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	q := r.URL.Query()
	count := validateQueryInt(q.Get("count"), 10)
	page := validateQueryInt(q.Get("page"), 1)

	id, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		log.Error(service.ParseErr, err)
		service.BadRequestResponse(w, service.ParseErr, err)
		return
	}

	tasks := Task{TodoListId: id}
	read, err := tasks.Read(count, page)
	if err != nil {
		log.Error(service.TaskReadErr, err)
		service.InternalServerErrorResponse(w, service.TaskReadErr, err)
		return
	}

	service.OkResponse(w, read)

	log.Info(service.TaskReadSuccess)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	listId, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		log.Error(service.ParseErr, err)
		service.BadRequestResponse(w, service.ParseErr, err)
		return
	}

	taskId, err := uuid.Parse(r.PathValue("taskId"))
	if err != nil {
		log.Error(service.ParseErr, err)
		service.BadRequestResponse(w, service.ParseErr, err)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(service.BodyReadErr, err)
		service.BadRequestResponse(w, service.BodyReadErr, err)
		return
	}

	task := Task{TaskId: taskId, TodoListId: listId}

	err = service.DeserializeJSON(data, &task)
	if err != nil {
		log.Error(service.JSONDeserializingErr, err)
		service.UnprocessableEntity(w, service.JSONDeserializingErr, err)
		return
	}

	err = task.validateTitle()
	if err != nil {
		log.Error(service.ValidationErr, err)
		service.BadRequestResponse(w, service.ValidationErr, err)
		return
	}

	err = task.Update()
	if err != nil {
		log.Error(service.TaskUpdateErr, err)
		service.InternalServerErrorResponse(w, service.TaskUpdateErr, err)
		return
	}

	service.OkResponse(w, task)

	log.WithFields(log.Fields{
		"id": task.ID,
	}).Info(service.TaskUpdateSuccess)
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	listId, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		log.Error(service.ParseErr, err)
		service.BadRequestResponse(w, service.ParseErr, err)
		return
	}

	taskId, err := uuid.Parse(r.PathValue("taskId"))
	if err != nil {
		log.Error(service.ParseErr, err)
		service.BadRequestResponse(w, service.ParseErr, err)
		return
	}

	t := Task{TodoListId: listId, TaskId: taskId}
	err = t.Delete()
	if err != nil {
		log.Error(service.TaskDeleteErr, err)
		service.InternalServerErrorResponse(w, service.TaskDeleteErr, err)
		return
	}

	service.OkResponse(w, service.DefaultResponse{
		ResultCode: 0,
		HttpCode:   http.StatusOK,
		Messages:   "",
		Data:       "",
	})

	log.WithFields(log.Fields{
		"id": t.ID,
	}).Info(service.TaskDeleteSuccess)
}
