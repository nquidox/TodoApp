package todoList

import (
	"github.com/google/uuid"
	"io"
	"net/http"
	"todoApp/api/service"
)

func CreateListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	title := RequestTitle{}
	todoList := &TodoList{}

	responseList := &TodoList{}
	item := Item{*responseList}
	response := Response{
		ResultCode: 0,
		Messages:   "",
		Data:       item,
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		serverResponse(w, ErrorResponse{
			ResultCode: http.StatusInternalServerError,
			Messages:   err.Error(),
			Data:       "",
		})
		return
	}

	err = service.DeserializeJSON(data, &title)
	if err != nil {
		serverResponse(w, ErrorResponse{
			ResultCode: http.StatusInternalServerError,
			Messages:   err.Error(),
			Data:       "",
		})
		return
	}

	err = validateListOnCreate(title.Title)
	if err != nil {
		serverResponse(w, ErrorResponse{
			ResultCode: http.StatusBadRequest,
			Messages:   err.Error(),
			Data:       "",
		})
		return
	}

	responseList, err = todoList.Create(title.Title)
	if err != nil {
		serverResponse(w, ErrorResponse{
			ResultCode: http.StatusInternalServerError,
			Messages:   err.Error(),
			Data:       "",
		})
		return
	}

	response.Data.List = *responseList
	serverResponse(w, response)
}

func GetAllListsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	todoLists := TodoList{}
	lists, err := todoLists.GetAllLists()
	if err != nil {
		serverResponse(w, ErrorResponse{
			ResultCode: http.StatusInternalServerError,
			Messages:   "Error getting lists: " + err.Error(),
			Data:       "",
		})
	}
	serverResponse(w, lists)
}

func UpdateListHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("listId"))
	todoList := TodoList{}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		serverResponse(w, ErrorResponse{
			ResultCode: http.StatusBadRequest,
			Messages:   err.Error(),
			Data:       "",
		})
		return
	}

	err = service.DeserializeJSON(data, &todoList)
	if err != nil {
		serverResponse(w, ErrorResponse{
			ResultCode: http.StatusBadRequest,
			Messages:   err.Error(),
			Data:       "",
		})
		return
	}

	err = todoList.Update(id, todoList.Title)
	if err != nil {
		serverResponse(w, ErrorResponse{
			ResultCode: http.StatusInternalServerError,
			Messages:   "Error updating list: " + err.Error(),
			Data:       "",
		})
		return
	}
}

func DeleteListHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		serverResponse(w, ErrorResponse{
			ResultCode: http.StatusBadRequest,
			Messages:   "ID Error: " + err.Error(),
			Data:       "",
		})
		return
	}

	todoList := TodoList{}
	err = todoList.Delete(id)
	if err != nil {
		serverResponse(w, ErrorResponse{
			ResultCode: http.StatusInternalServerError,
			Messages:   "Internal server error: " + err.Error(),
			Data:       "",
		})
		return
	}
}
