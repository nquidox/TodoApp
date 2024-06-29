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
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusBadRequest,
			Messages:   "Error reading body: " + err.Error(),
			Data:       "",
		})
		return
	}

	err = service.DeserializeJSON(data, &todoList)
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusInternalServerError,
			Messages:   "JSON read error: " + err.Error(),
			Data:       "",
		})
		return
	}

	err = todoList.validateTitle()
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusBadRequest,
			Messages:   "Validation error: " + err.Error(),
			Data:       "",
		})
		return
	}

	err = todoList.Create()
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusInternalServerError,
			Messages:   "Create list error: " + err.Error(),
			Data:       "",
		})
		return
	}

	service.ServerResponse(w, service.ErrorResponse{
		ResultCode: 0,
		ErrorCode:  http.StatusOK,
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
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusInternalServerError,
			Messages:   "Error getting lists: " + err.Error(),
			Data:       "",
		})
	}
	service.ServerResponse(w, lists)
}

func UpdateListHandler(w http.ResponseWriter, r *http.Request) {
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

	todoList := TodoList{Uuid: id}

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

	err = service.DeserializeJSON(data, &todoList)
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusBadRequest,
			Messages:   "JSON read error: " + err.Error(),
			Data:       "",
		})
		return
	}

	err = todoList.Update()
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusInternalServerError,
			Messages:   "Error updating list: " + err.Error(),
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

func DeleteListHandler(w http.ResponseWriter, r *http.Request) {
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
	todoList := TodoList{Uuid: id}

	err = todoList.Delete()
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusInternalServerError,
			Messages:   "Internal server error: " + err.Error(),
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
