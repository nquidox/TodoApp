package todoList

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"todoApp/api/service"
	"todoApp/api/user"
)

// CreateListHandler godoc
//
//	@Summary		Create todo list
//	@Description	Creates new todo list
//	@Tags			Todo lists
//	@Security		BasicAuth
//	@Accept			json
//	@Produce		json
//	@Param			model	body		createTodoList			true	"Create new todo list"
//	@Success		200		{object}	service.DefaultResponse	"OK"
//	@Failure		400		{object}	service.errorResponse	"Bad request"
//	@Failure		401		{object}	service.errorResponse	"Unauthorized"
//	@Failure		500		{object}	service.errorResponse	"Internal server error"
//	@Router			/todo-lists [post]
func CreateListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	todoList := createTodoList{}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(service.BodyReadErr, err)
		service.BadRequestResponse(w, service.BodyReadErr, err)
		return
	}

	err = service.DeserializeJSON(data, &todoList)
	if err != nil {
		log.Error(service.JSONDeserializingErr, err)
		service.UnprocessableEntityResponse(w, service.JSONReadErr, err)
		return
	}

	err = todoList.validateTitle()
	if err != nil {
		log.Error(service.ValidationErr, err)
		service.BadRequestResponse(w, service.ValidationErr, err)
		return
	}

	todoList.ListUuid = uuid.New()
	err = todoList.Create(dbWorker)
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
		"id":    todoList.ListUuid,
		"title": todoList.Title,
	}).Info(service.TodoListCreateSuccess)
}

// GetAllListsHandler godoc
//
//	@Summary		Get todo lists
//	@Description	Requests all todo list
//	@Tags			Todo lists
//	@Security		BasicAuth
//	@Produce		json
//	@Success		200	{array}		readTodoList			"OK"
//	@Success		204	{array}		readTodoList			"No Content"
//	@Failure		400	{object}	service.errorResponse	"Bad request"
//	@Failure		401	{object}	service.errorResponse	"Unauthorized"
//	@Failure		404	{object}	service.errorResponse	"Not Found"
//	@Failure		500	{object}	service.errorResponse	"Internal server error"
//	@Router			/todo-lists [get]
func GetAllListsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	todoLists := readTodoList{}

	//FIXME: refactor this in the future
	token, err := r.Cookie("token")
	if err != nil {
		log.Error("Temporary code!!! Cookie error: ", err)
	}
	log.Debug("token recieved: ", token.Value)

	s := user.Session{Token: token.Value}

	err = s.Read(dbWorker)
	if err != nil {
		log.Error("Temporary code!!! Session read error: ", err)
	}
	log.Debug("session recieved: ", s)

	authUsr, err := authWorker.IsUserLoggedIn(dbWorker, s.UserUuid)
	if err != nil {
		log.Error("Temporary code!!! authUser interface method error: ", err)
	}
	log.Debug("authUser recieved: ", authUsr)

	//for now we ignore IsSuperuser
	lists, err := todoLists.GetAllLists(dbWorker, authUsr)
	if err != nil {
		if err.Error() == "404" {
			service.OkResponse(w, service.DefaultResponse{
				ResultCode: 0,
				HttpCode:   http.StatusNoContent,
				Messages:   "",
				Data:       nil,
			})
			log.Error(service.DBNotFound)
			return
		}
		log.Error(service.ListReadErr, err)
		service.InternalServerErrorResponse(w, service.ListReadErr, err)
		return
	}

	service.OkResponse(w, lists)
	log.Info(service.TodoListReadSuccess)
}

// UpdateListHandler godoc
//
//	@Summary		Update todo list
//	@Description	Updates todo list
//	@Tags			Todo lists
//	@Security		BasicAuth
//	@Produce		json
//	@Param			listId	path		string					true	"List uuid"
//	@Param			data	body		createTodoList			true	"List data for update"
//	@Success		200		{object}	TodoList				"OK"
//	@Failure		400		{object}	service.errorResponse	"Bad request"
//	@Failure		401		{object}	service.errorResponse	"Unauthorized"
//	@Failure		404		{object}	service.errorResponse	"Not Found"
//	@Failure		500		{object}	service.errorResponse	"Internal server error"
//	@Router			/todo-lists/{listId} [put]
func UpdateListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		log.Error(service.ParseErr, err)
		service.BadRequestResponse(w, service.ParseErr, err)
		return
	}

	todoList := createTodoList{ListUuid: id}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(service.BodyReadErr, err)
		service.BadRequestResponse(w, service.BodyReadErr, err)
		return
	}

	err = service.DeserializeJSON(data, &todoList)
	if err != nil {
		log.Error(service.JSONDeserializingErr, err)
		service.UnprocessableEntityResponse(w, service.JSONReadErr, err)
		return
	}

	err = todoList.Update(dbWorker)
	if err != nil {
		if err.Error() == "404" {
			service.NotFoundResponse(w, "")
			log.Error(service.DBNotFound)
			return
		}
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
		"id": todoList.ListUuid,
	}).Info(service.TodoListUpdateSuccess)
}

// DeleteListHandler godoc
//
//	@Summary		Delete todo list
//	@Description	Deletes todo list
//	@Tags			Todo lists
//	@Security		BasicAuth
//	@Produce		json
//	@Param			listId	path		string					true	"list uuid"
//	@Success		200		{object}	service.DefaultResponse	"OK"
//	@Failure		400		{object}	service.errorResponse	"Bad request"
//	@Failure		401		{object}	service.errorResponse	"Unauthorized"
//	@Failure		404		{object}	service.errorResponse	"Not Found"
//	@Failure		500		{object}	service.errorResponse	"Internal server error"
//	@Router			/todo-lists/{listId} [delete]
func DeleteListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := uuid.Parse(r.PathValue("listId"))
	if err != nil {
		log.Error(service.ParseErr, err)
		service.BadRequestResponse(w, service.ParseErr, err)
		return
	}
	todoList := TodoList{ListUuid: id}

	err = todoList.Delete(dbWorker)

	if err != nil {
		if err.Error() == "404" {
			service.NotFoundResponse(w, "")
			log.Error(service.DBNotFound)
			return
		}
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
		"id": todoList.ListUuid,
	}).Info(service.TodoListDeleteSuccess)
}
