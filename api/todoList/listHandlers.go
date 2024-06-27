package todoList

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"todoApp/api/service"
)

type ListTitle struct {
	Title string `json:"title"`
}

func CreateListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	title := ListTitle{}
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
		response.ResultCode = http.StatusInternalServerError
		response.Messages = err.Error()
		serverResponse(w, response)
		return
	}

	err = service.DeserializeJSON(data, &title)
	if err != nil {
		response.ResultCode = http.StatusInternalServerError
		response.Messages = err.Error()
		serverResponse(w, response)
		return
	}

	err = validateListOnCreate(title.Title)
	if err != nil {
		response.ResultCode = http.StatusBadRequest
		response.Messages = "Error: " + err.Error()
		serverResponse(w, response)
		return
	}

	responseList, err = todoList.Create(title.Title)
	if err != nil {
		response.ResultCode = http.StatusInternalServerError
		response.Messages = err.Error()
		serverResponse(w, response)
		return
	}

	response.Data.List = *responseList
	serverResponse(w, response)
}

func GetAllListsHandler(w http.ResponseWriter, r *http.Request) {
	todoLists := TodoList{}
	lists, err := todoLists.GetAllLists()
	if err != nil {
		fmt.Fprint(w, "Error getting lists: ", err, http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		response, _ := service.SerializeJSON(lists)
		fmt.Fprint(w, string(response))
	}
}

func UpdateListHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("listId"))
	todoList := TodoList{}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error parsing body", http.StatusBadRequest)
	}

	err = service.DeserializeJSON(data, &todoList)
	if err != nil {
		http.Error(w, "Error parsing body", http.StatusBadRequest)
	}

	err = todoList.Update(id, todoList.Title)
	if err != nil {
		fmt.Fprint(w, "Error updating list: ", err, http.StatusInternalServerError)
	}
}

func DeleteListHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		response := Response{
			ResultCode: http.StatusBadRequest,
			Messages:   "ID Error: " + err.Error(),
			Data:       Item{},
		}
		serverResponse(w, response)
		return
	}

	todoList := TodoList{}
	err = todoList.Delete(id)
	if err != nil {
		response := Response{
			ResultCode: 500,
			Messages:   "Error deleting list: " + err.Error(),
			Data:       Item{},
		}
		serverResponse(w, response)
	}
}

func serverResponse(w http.ResponseWriter, response Response) {
	w.Header().Set("Content-Type", "application/json")
	bytes, err := service.SerializeJSON(response)
	if err != nil {
		http.Error(w, "Error serializing response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, "Error writing response: "+err.Error(), http.StatusInternalServerError)
	}
}
