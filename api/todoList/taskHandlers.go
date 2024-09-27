package todoList

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"todoApp/api/service"
)

// createTaskFunc godoc
//
//	@Security		BasicAuth
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
//	@Failure		422		{object}	service.errorResponse	"Unprocessable entity"
//	@Failure		500		{object}	service.errorResponse	"Internal server error"
//	@Router			/todo-lists/{listId}/tasks [post]
func createTaskFunc(s *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var aUser authUser
		err := aUser.isAuth(w, r, s)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Error(service.Unauthorized, err)
			service.UnauthorizedResponse(w, "")
			return
		}

		id, err := uuid.Parse(r.PathValue("listId"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(service.ParseErr, err)
			service.BadRequestResponse(w, service.ParseErr, err)
			return
		}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(service.BodyReadErr, err)
			service.BadRequestResponse(w, service.BodyReadErr, err)
			return
		}

		task := createTask{TodoListUUID: id}

		err = service.DeserializeJSON(data, &task)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			log.Error(service.JSONDeserializingErr, err)
			service.UnprocessableEntityResponse(w, service.JSONDeserializingErr, err)
			return
		}

		err = task.validateTitle()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(service.ValidationErr, err)
			service.BadRequestResponse(w, service.ValidationErr, err)
			return
		}

		newTask := Task{
			Description:  task.Description,
			Title:        task.Title,
			Completed:    "",
			Status:       task.Status,
			Priority:     task.Priority,
			StartDate:    task.StartDate,
			Deadline:     task.Deadline,
			TaskUUID:     uuid.New(),
			TodoListUUID: task.TodoListUUID,
			Order:        task.Order,
			OwnerUUID:    aUser.UserUUID,
		}

		err = newTask.Create(s.DbWorker)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(service.TaskCreateErr, err)
			service.InternalServerErrorResponse(w, service.TaskCreateErr, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.WithFields(log.Fields{
			"Title": task.Title,
		}).Info(service.TaskCreateSuccess)

		service.OkResponse(w, task)
	}
}

// getTaskFunc godoc
//
//	@Security		BasicAuth
//	@Summary		Get tasks
//	@Description	Requests all tasks with query parameters. Order, count and page params are optional. Defaults: order= desc, count=10, page=1
//	@Tags			Tasks
//	@Produce		json
//	@Param			listId	path		string					true	"list uuid"
//	@Param			order	query		string					false	"asc/desc (default)"
//	@Param			count	query		string					false	"Count (number of task to show per page)"
//	@Param			page	query		string					false	"Page number"
//	@Success		200		{array}		Task					"OK"
//	@Success		204		{object}	service.DefaultResponse	"No Content"
//	@Failure		400		{object}	service.errorResponse	"Bad request"
//	@Failure		401		{object}	service.errorResponse	"Unauthorized"
//	@Failure		500		{object}	service.errorResponse	"Internal server error"
//	@Router			/todo-lists/{listId}/tasks [get]
func getTaskFunc(s *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var aUser authUser
		err := aUser.isAuth(w, r, s)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Error(service.Unauthorized, err)
			service.UnauthorizedResponse(w, "")
			return
		}

		q := r.URL.Query()
		order := validateOrder(r.URL.Query().Get("order"))
		count := validateQueryInt(q.Get("count"), 10)
		page := validateQueryInt(q.Get("page"), 1)

		id, err := uuid.Parse(r.PathValue("listId"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(service.ParseErr, err)
			service.BadRequestResponse(w, service.ParseErr, err)
			return
		}

		tasks := Task{TodoListUUID: id, OwnerUUID: aUser.UserUUID}
		log.WithFields(log.Fields{
			"ListId": id,
			"Order":  order,
			"Count":  count,
			"Page":   page,
		}).Debug("Query and path params")

		read, err := tasks.Read(s.DbWorker, order, count, page)
		if err != nil {
			if err.Error() == "404" {
				w.WriteHeader(http.StatusNoContent)
				log.Info(service.NoContent)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(service.TaskReadErr, err)
			service.InternalServerErrorResponse(w, service.TaskReadErr, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info(service.TaskReadSuccess)
		service.OkResponse(w, read)

	}
}

// updateTaskFunc godoc
//
//	@Security		BasicAuth
//	@Summary		Update task
//	@Description	Updates task
//	@Tags			Tasks
//	@Produce		json
//	@Param			data	body		createTask				true	"Task data for update"
//	@Success		200		{object}	createTask				"OK"
//	@Failure		400		{object}	service.errorResponse	"Bad request"
//	@Failure		401		{object}	service.errorResponse	"Unauthorized"
//	@Failure		404		{object}	service.errorResponse	"Not Found"
//	@Failure		422		{object}	service.errorResponse	"Unprocessable entity"
//	@Failure		500		{object}	service.errorResponse	"Internal server error"
//	@Router			/todo-lists/{listId}/tasks/{taskId} [put]
func updateTaskFunc(s *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var aUser authUser
		err := aUser.isAuth(w, r, s)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Error(service.Unauthorized, err)
			service.UnauthorizedResponse(w, "")
			return
		}

		listId, err := uuid.Parse(r.PathValue("listId"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(service.ParseErr, err)
			service.BadRequestResponse(w, service.ParseErr, err)
			return
		}

		taskId, err := uuid.Parse(r.PathValue("taskId"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(service.ParseErr, err)
			service.BadRequestResponse(w, service.ParseErr, err)
			return
		}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(service.BodyReadErr, err)
			service.BadRequestResponse(w, service.BodyReadErr, err)
			return
		}

		task := createTask{TaskUUID: taskId, TodoListUUID: listId, OwnerUUID: aUser.UserUUID}

		err = service.DeserializeJSON(data, &task)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			log.Error(service.JSONDeserializingErr, err)
			service.UnprocessableEntityResponse(w, service.JSONDeserializingErr, err)
			return
		}

		err = task.validateTitle()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(service.ValidationErr, err)
			service.BadRequestResponse(w, service.ValidationErr, err)
			return
		}

		err = task.Update(s.DbWorker)
		if err != nil {
			if err.Error() == "404" {
				w.WriteHeader(http.StatusNotFound)
				log.Error(service.DBNotFound)
				service.NotFoundResponse(w, "")
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(service.TaskUpdateErr, err)
			service.InternalServerErrorResponse(w, service.TaskUpdateErr, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.WithFields(log.Fields{
			"task id": task.TaskUUID,
		}).Info(service.TaskUpdateSuccess)
		service.OkResponse(w, task)
	}
}

// deleteTaskFunc godoc
//
//	@Security		BasicAuth
//	@Summary		Delete task
//	@Description	Deletes task
//	@Tags			Tasks
//	@Produce		json
//	@Success		200	{object}	service.DefaultResponse	"OK"
//	@Failure		400	{object}	service.errorResponse	"Bad request"
//	@Failure		401	{object}	service.errorResponse	"Unauthorized"
//	@Failure		404	{object}	service.errorResponse	"Not Found"
//	@Failure		500	{object}	service.errorResponse	"Internal server error"
//	@Router			/todo-lists/{listId}/tasks/{taskId} [delete]
func deleteTaskFunc(s *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var aUser authUser
		err := aUser.isAuth(w, r, s)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Error(service.Unauthorized, err)
			service.UnauthorizedResponse(w, "")
			return
		}

		listId, err := uuid.Parse(r.PathValue("listId"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(service.ParseErr, err)
			service.BadRequestResponse(w, service.ParseErr, err)
			return
		}

		taskId, err := uuid.Parse(r.PathValue("taskId"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(service.ParseErr, err)
			service.BadRequestResponse(w, service.ParseErr, err)
			return
		}

		t := Task{TodoListUUID: listId, TaskUUID: taskId, OwnerUUID: aUser.UserUUID}
		err = t.Delete(s.DbWorker)
		if err != nil {
			if err.Error() == "404" {
				w.WriteHeader(http.StatusNotFound)
				log.Error(service.DBNotFound)
				service.NotFoundResponse(w, "")
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(service.TaskDeleteErr, err)
			service.InternalServerErrorResponse(w, service.TaskDeleteErr, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.WithFields(log.Fields{
			"id": t.ID,
		}).Info(service.TaskDeleteSuccess)

		service.OkResponse(w, service.DefaultResponse{
			ResultCode: 0,
			HttpCode:   http.StatusOK,
			Messages:   "",
			Data:       "",
		})
	}
}
