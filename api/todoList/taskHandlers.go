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
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error parsing body", http.StatusBadRequest)
	}

	title := ""
	err = service.DeserializeJSON(data, title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	task := Task{}

	err = task.Create(id, title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	bytes, err := service.SerializeJSON(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	q := r.URL.Query()
	count := validateQueryInt(q.Get("count"), 10)
	page := validateQueryInt(q.Get("page"), 1)

	id, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	tasks := Task{}
	read, err := tasks.Read(id, count, page)
	if err != nil {
		return
	}

	bytes, err := service.SerializeJSON(read)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	listId, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	taskId, err := uuid.Parse(r.PathValue("taskId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error parsing body", http.StatusBadRequest)
	}

	task := Task{}

	err = service.DeserializeJSON(data, &task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	result, err := task.Update(listId, taskId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	bytes, err := service.SerializeJSON(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	listId, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	taskId, err := uuid.Parse(r.PathValue("taskId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
