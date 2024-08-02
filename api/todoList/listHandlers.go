package todoList

import (
	"github.com/google/uuid"
	"io"
	"net/http"
	"todoApp/api/service"
)

func CreateListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	todoList := TodoList{}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		service.BadRequestResponse(w, service.BodyReadErr, err)
		return
	}

	err = service.DeserializeJSON(data, &todoList)
	if err != nil {
		service.UnprocessableEntity(w, service.JSONReadErr, err)
		return
	}

	err = todoList.validateTitle()
	if err != nil {
		service.BadRequestResponse(w, service.ValidationErr, err)
		return
	}

	err = todoList.Create()
	if err != nil {
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
}

func GetAllListsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	todoLists := TodoList{}

	lists, err := todoLists.GetAllLists()
	if err != nil {
		service.InternalServerErrorResponse(w, service.ListReadErr, err)
		return
	}
	service.OkResponse(w, lists)
}

func UpdateListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		service.BadRequestResponse(w, service.ParseErr, err)
		return
	}

	todoList := TodoList{Uuid: id}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		service.BadRequestResponse(w, service.BodyReadErr, err)
		return
	}

	err = service.DeserializeJSON(data, &todoList)
	if err != nil {
		service.UnprocessableEntity(w, service.JSONReadErr, err)
		return
	}

	err = todoList.Update()
	if err != nil {
		service.InternalServerErrorResponse(w, service.ListUpdateErr, err)
		return
	}

	service.OkResponse(w, service.DefaultResponse{
		ResultCode: 0,
		HttpCode:   http.StatusOK,
		Messages:   "",
		Data:       "",
	})
}

func DeleteListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		service.BadRequestResponse(w, service.ParseErr, err)
		return
	}
	todoList := TodoList{Uuid: id}

	err = todoList.Delete()
	if err != nil {
		service.InternalServerErrorResponse(w, service.ListDeleteErr, err)
		return
	}

	service.OkResponse(w, service.DefaultResponse{
		ResultCode: 0,
		HttpCode:   http.StatusOK,
		Messages:   "",
		Data:       "",
	})
}
