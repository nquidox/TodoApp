package todoList

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"todoApp/api/service"
)

// createListFunc godoc
//
//	@Security		BasicAuth
//	@Summary		Create todo list
//	@Description	Creates new todo list
//	@Tags			Todo lists
//	@Accept			json
//	@Produce		json
//	@Param			model	body		createTodoList			true	"Create new todo list"
//	@Success		200		{object}	service.DefaultResponse	"OK"
//	@Failure		400		{object}	service.errorResponse	"Bad request"
//	@Failure		401		{object}	service.errorResponse	"Unauthorized"
//	@Failure		422		{object}	service.errorResponse	"Unprocessable entity"
//	@Failure		500		{object}	service.errorResponse	"Internal server error"
//	@Router			/todo-lists [post]
func createListFunc(s *Service) http.HandlerFunc {
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

		todoList := createTodoList{}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(service.BodyReadErr, err)
			service.BadRequestResponse(w, service.BodyReadErr, err)
			return
		}

		err = service.DeserializeJSON(data, &todoList)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			log.Error(service.JSONDeserializingErr, err)
			service.UnprocessableEntityResponse(w, service.JSONReadErr, err)
			return
		}

		err = todoList.validateTitle()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(service.ValidationErr, err)
			service.BadRequestResponse(w, service.ValidationErr, err)
			return
		}

		todoList.ListUuid = uuid.New()
		todoList.OwnerUuid = aUser.UserUUID
		err = todoList.Create(s.DbWorker)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(service.ListCreateErr, err)
			service.InternalServerErrorResponse(w, service.ListCreateErr, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.WithFields(log.Fields{
			"id":    todoList.ListUuid,
			"title": todoList.Title,
		}).Info(service.TodoListCreateSuccess)

		service.OkResponse(w, service.DefaultResponse{
			ResultCode: 0,
			HttpCode:   http.StatusOK,
			Messages:   "",
			Data: Item{
				List: todoList,
			},
		})
	}
}

// getAllListsFunc godoc
//
//	@Security		BasicAuth
//	@Summary		Get todo lists
//	@Description	Requests all todo list
//	@Tags			Todo lists
//	@Produce		json
//	@Success		200	{array}		readTodoList			"OK"
//	@Success		204	{array}		service.DefaultResponse	"No Content"
//	@Failure		401	{object}	service.errorResponse	"Unauthorized"
//	@Failure		500	{object}	service.errorResponse	"Internal server error"
//	@Router			/todo-lists [get]
func getAllListsFunc(s *Service) http.HandlerFunc {
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

		todoLists := readTodoList{}
		lists, err := todoLists.GetAllLists(s.DbWorker, aUser)
		if err != nil {
			if err.Error() == "404" {
				w.WriteHeader(http.StatusNoContent)
				log.Info(service.NoContent)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(service.ListReadErr, err)
			service.InternalServerErrorResponse(w, service.ListReadErr, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info(service.TodoListReadSuccess)
		service.OkResponse(w, lists)
	}
}

// updateListFunc godoc
//
//	@Security		BasicAuth
//	@Summary		Update todo list
//	@Description	Updates todo list
//	@Tags			Todo lists
//	@Produce		json
//	@Param			data	body		createTodoList			true	"List data for update"
//	@Success		200		{object}	TodoList				"OK"
//	@Failure		400		{object}	service.errorResponse	"Bad request"
//	@Failure		401		{object}	service.errorResponse	"Unauthorized"
//	@Failure		404		{object}	service.errorResponse	"Not Found"
//	@Failure		422		{object}	service.errorResponse	"Unprocessable entity"
//	@Failure		500		{object}	service.errorResponse	"Internal server error"
//	@Router			/todo-lists/{listId} [put]
func updateListFunc(s *Service) http.HandlerFunc {
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

		todoList := createTodoList{ListUuid: id, OwnerUuid: aUser.UserUUID}
		err = service.DeserializeJSON(data, &todoList)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			log.Error(service.JSONDeserializingErr, err)
			service.UnprocessableEntityResponse(w, service.JSONReadErr, err)
			return
		}

		err = todoList.Update(s.DbWorker)
		if err != nil {
			if err.Error() == "404" {
				w.WriteHeader(http.StatusNotFound)
				log.Error(service.DBNotFound)
				service.NotFoundResponse(w, "")
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(service.ListUpdateErr, err)
			service.InternalServerErrorResponse(w, service.ListUpdateErr, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.WithFields(log.Fields{
			"id": todoList.ListUuid,
		}).Info(service.TodoListUpdateSuccess)

		service.OkResponse(w, service.DefaultResponse{
			ResultCode: 0,
			HttpCode:   http.StatusOK,
			Messages:   "",
			Data:       "",
		})
	}
}

// deleteListFunc godoc
//
//	@Security		BasicAuth
//	@Summary		Delete todo list
//	@Description	Deletes todo list
//	@Tags			Todo lists
//	@Produce		json
//	@Success		200	{object}	service.DefaultResponse	"OK"
//	@Failure		400	{object}	service.errorResponse	"Bad request"
//	@Failure		401	{object}	service.errorResponse	"Unauthorized"
//	@Failure		404	{object}	service.errorResponse	"Not Found"
//	@Failure		500	{object}	service.errorResponse	"Internal server error"
//	@Router			/todo-lists/{listId} [delete]
func deleteListFunc(s *Service) http.HandlerFunc {
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

		todoList := TodoList{ListUuid: id, OwnerUuid: aUser.UserUUID}
		err = todoList.Delete(s.DbWorker)

		if err != nil {
			if err.Error() == "404" {
				w.WriteHeader(http.StatusNotFound)
				log.Error(service.DBNotFound)
				service.NotFoundResponse(w, "")
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(service.ListDeleteErr, err)
			service.InternalServerErrorResponse(w, service.ListDeleteErr, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.WithFields(log.Fields{
			"id": todoList.ListUuid,
		}).Info(service.TodoListDeleteSuccess)

		service.OkResponse(w, service.DefaultResponse{
			ResultCode: 0,
			HttpCode:   http.StatusOK,
			Messages:   "",
			Data:       "",
		})
	}
}
