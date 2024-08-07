package todoList

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"todoApp/api/service"
)

func CreateListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	todoList := TodoList{}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(service.BodyReadErr, err)
		service.BadRequestResponse(w, service.BodyReadErr, err)
		return
	}

	err = service.DeserializeJSON(data, &todoList)
	if err != nil {
		log.Error(service.JSONDeserializingErr, err)
		service.UnprocessableEntity(w, service.JSONReadErr, err)
		return
	}

	err = todoList.validateTitle()
	if err != nil {
		log.Error(service.ValidationErr, err)
		service.BadRequestResponse(w, service.ValidationErr, err)
		return
	}

	err = todoList.Create()
	if err != nil {
		log.Error(service.ListCreateErr, err)
		service.InternalServerErrorResponse(w, service.ListCreateErr, err)
		return
	}

	service.OkResponse(w, service.DefaultResponse{
		ResultCode: 0,
		HttpCode:   http.StatusOK,
		Messages:   "",
		Data: Item{
			List: todoList,
		},
	})

	log.WithFields(log.Fields{
		"id":    todoList.Uuid,
		"title": todoList.Title,
	}).Info(service.TodoListCreateSuccess)
}

func GetAllListsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	todoLists := TodoList{}

	lists, err := todoLists.GetAllLists()
	if err != nil {
		log.Error(service.ListReadErr, err)
		service.InternalServerErrorResponse(w, service.ListReadErr, err)
		return
	}

	service.OkResponse(w, lists)
	log.Info(service.TodoListReadSuccess)
}

func UpdateListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		log.Error(service.ParseErr, err)
		service.BadRequestResponse(w, service.ParseErr, err)
		return
	}

	todoList := TodoList{Uuid: id}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(service.BodyReadErr, err)
		service.BadRequestResponse(w, service.BodyReadErr, err)
		return
	}

	err = service.DeserializeJSON(data, &todoList)
	if err != nil {
		log.Error(service.JSONDeserializingErr, err)
		service.UnprocessableEntity(w, service.JSONReadErr, err)
		return
	}

	err = todoList.Update()
	if err != nil {
		log.Error(service.ListUpdateErr, err)
		service.InternalServerErrorResponse(w, service.ListUpdateErr, err)
		return
	}

	service.OkResponse(w, service.DefaultResponse{
		ResultCode: 0,
		HttpCode:   http.StatusOK,
		Messages:   "",
		Data:       "",
	})

	log.WithFields(log.Fields{
		"id": todoList.Uuid,
	}).Info(service.TodoListUpdateSuccess)
}

func DeleteListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		log.Error(service.ParseErr, err)
		service.BadRequestResponse(w, service.ParseErr, err)
		return
	}
	todoList := TodoList{Uuid: id}

	err = todoList.Delete()
	if err != nil {
		log.Error(service.ListDeleteErr, err)
		service.InternalServerErrorResponse(w, service.ListDeleteErr, err)
		return
	}

	service.OkResponse(w, service.DefaultResponse{
		ResultCode: 0,
		HttpCode:   http.StatusOK,
		Messages:   "",
		Data:       "",
	})

	log.WithFields(log.Fields{
		"id": todoList.Uuid,
	}).Info(service.TodoListDeleteSuccess)
}
