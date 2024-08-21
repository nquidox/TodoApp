package todoList

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"todoApp/api/service"
)

// CreateTaskHandler godoc
//
//	@Summary		Create task list
//	@Description	Creates new task. Time format example: "02-01-2006 15:04:05"
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Param			listId	path		string					true	"List UUID"
//	@Param			model	body		createTask				true	"Create new task"
//	@Success		200		{object}	createTask				"OK"
//	@Failure		400		{object}	service.errorResponse	"Bad request"
//	@Failure		401		{object}	service.errorResponse	"Unauthorized"
//	@Failure		500		{object}	service.errorResponse	"Internal server error"
//	@Router			/todo-lists/{listId}/tasks [post]
func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		service.BadRequestResponse(w, service.ParseErr, err)
		log.Error(service.ParseErr, err)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		service.BadRequestResponse(w, service.BodyReadErr, err)
		log.Error(service.BodyReadErr, err)
		return
	}

	task := createTask{TodoListId: id}

	err = service.DeserializeJSON(data, &task)
	if err != nil {
		service.UnprocessableEntityResponse(w, service.JSONDeserializingErr, err)
		log.Error(service.JSONDeserializingErr, err)
		return
	}

	err = task.validateTitle()
	if err != nil {
		service.BadRequestResponse(w, service.ValidationErr, err)
		log.Error(service.ValidationErr, err)
		return
	}

	newTask := Task{
		Description: task.Description,
		Title:       task.Title,
		Completed:   "",
		Status:      task.Status,
		Priority:    task.Priority,
		StartDate:   validateTime(task.StartDate),
		Deadline:    validateTime(task.Deadline),
		TaskId:      uuid.New(),
		TodoListId:  task.TodoListId,
		Order:       task.Order,
	}

	err = newTask.Create(dbw)
	if err != nil {
		service.InternalServerErrorResponse(w, service.TaskCreateErr, err)
		log.Error(service.TaskCreateErr, err)
		return
	}

	service.OkResponse(w, task)

	log.WithFields(log.Fields{
		"Title": task.Title,
	}).Info(service.TaskCreateSuccess)
}

// GetTaskHandler godoc
//
//	@Summary		Get tasks
//	@Description	Requests all tasks with query parameters. Count and page params are optional. Defaults: count=10, page=1
//	@Tags			Tasks
//	@Produce		json
//	@Param			listId	path		string					true	"list uuid"
//	@Param			count	query		string					false	"Count (number of task to show per page)"
//	@Param			page	query		string					false	"Page number"
//	@Success		200		{array}		Task					"OK"
//	@Success		204		{object}	service.DefaultResponse	"No Content"
//	@Failure		400		{object}	service.errorResponse	"Bad request"
//	@Failure		401		{object}	service.errorResponse	"Unauthorized"
//	@Failure		500		{object}	service.errorResponse	"Internal server error"
//	@Router			/todo-lists/{listId}/tasks [get]
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	q := r.URL.Query()
	count := validateQueryInt(q.Get("count"), 10)
	page := validateQueryInt(q.Get("page"), 1)

	id, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		service.BadRequestResponse(w, service.ParseErr, err)
		log.Error(service.ParseErr, err)
		return
	}

	tasks := Task{TodoListId: id}
	log.WithFields(log.Fields{
		"ListId": id,
		"Count":  count,
		"Page":   page,
	}).Debug("Query and path params")

	read, err := tasks.Read(dbw, count, page)
	if err != nil {
		//if no records found, return success 204 no content instead of 404
		if err.Error() == "404" {
			service.OkResponse(w, service.DefaultResponse{
				ResultCode: 0,
				HttpCode:   http.StatusNoContent,
				Messages:   "",
				Data:       nil,
			})
			return
		}
		service.InternalServerErrorResponse(w, service.TaskReadErr, err)
		log.Error(service.TaskReadErr, err)
		return
	}

	service.OkResponse(w, read)

	log.Info(service.TaskReadSuccess)
}

// UpdateTaskHandler godoc
//
//	@Summary		Update task
//	@Description	Updates task
//	@Tags			Tasks
//	@Produce		json
//	@Param			listId	path		string					true	"list uuid"
//	@Param			taskId	path		string					true	"task uuid"
//	@Param			data	body		createTask				true	"Task data for update"
//	@Success		200		{object}	createTask				"OK"
//	@Failure		400		{object}	service.errorResponse	"Bad request"
//	@Failure		401		{object}	service.errorResponse	"Unauthorized"
//	@Failure		404		{object}	service.errorResponse	"Not Found"
//	@Failure		500		{object}	service.errorResponse	"Internal server error"
//	@Router			/todo-lists/{listId}/tasks/{taskId} [put]
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	listId, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		service.BadRequestResponse(w, service.ParseErr, err)
		log.Error(service.ParseErr, err)
		return
	}

	taskId, err := uuid.Parse(r.PathValue("taskId"))
	if err != nil {
		service.BadRequestResponse(w, service.ParseErr, err)
		log.Error(service.ParseErr, err)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		service.BadRequestResponse(w, service.BodyReadErr, err)
		log.Error(service.BodyReadErr, err)
		return
	}

	task := createTask{TaskId: taskId, TodoListId: listId}

	err = service.DeserializeJSON(data, &task)
	if err != nil {
		service.UnprocessableEntityResponse(w, service.JSONDeserializingErr, err)
		log.Error(service.JSONDeserializingErr, err)
		return
	}

	err = task.validateTitle()
	if err != nil {
		service.BadRequestResponse(w, service.ValidationErr, err)
		log.Error(service.ValidationErr, err)
		return
	}

	err = task.Update(dbw)
	if err != nil {
		service.InternalServerErrorResponse(w, service.TaskUpdateErr, err)
		log.Error(service.TaskUpdateErr, err)
		return
	}

	service.OkResponse(w, task)

	log.WithFields(log.Fields{
		"task id": task.TaskId,
	}).Info(service.TaskUpdateSuccess)
}

// DeleteTaskHandler godoc
//
//	@Summary		Delete task
//	@Description	Deletes task
//	@Tags			Tasks
//	@Produce		json
//	@Param			listId	path		string					true	"list uuid"
//	@Param			taskId	path		string					true	"task uuid"
//	@Success		200		{object}	service.DefaultResponse	"OK"
//	@Failure		400		{object}	service.errorResponse	"Bad request"
//	@Failure		401		{object}	service.errorResponse	"Unauthorized"
//	@Failure		404		{object}	service.errorResponse	"Not Found"
//	@Failure		500		{object}	service.errorResponse	"Internal server error"
//	@Router			/todo-lists/{listId}/tasks/{taskId} [delete]
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	listId, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		service.BadRequestResponse(w, service.ParseErr, err)
		log.Error(service.ParseErr, err)
		return
	}

	taskId, err := uuid.Parse(r.PathValue("taskId"))
	if err != nil {
		service.BadRequestResponse(w, service.ParseErr, err)
		log.Error(service.ParseErr, err)
		return
	}

	t := Task{TodoListId: listId, TaskId: taskId}
	err = t.Delete(dbw)
	if err != nil {
		service.InternalServerErrorResponse(w, service.TaskDeleteErr, err)
		log.Error(service.TaskDeleteErr, err)
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
