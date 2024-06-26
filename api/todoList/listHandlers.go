package todoList

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"todoApp/api/service"
)

func CreateListHandler(w http.ResponseWriter, r *http.Request) {
	todoList := TodoList{}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error parsing body", http.StatusBadRequest)
	}

	err = service.DeserializeJSON(data, &todoList)
	if err != nil {
		http.Error(w, "Error parsing body", http.StatusBadRequest)
	}

	err = todoList.Create(todoList.Title)
	if err != nil {
		fmt.Fprint(w, "Error creating list: ", err, http.StatusBadRequest)
	}
	fmt.Println(todoList)
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

	todoList := TodoList{}
	err = todoList.Delete(id)
	if err != nil {
		fmt.Fprint(w, "Error deleting list: ", err, http.StatusInternalServerError)
	}
}
