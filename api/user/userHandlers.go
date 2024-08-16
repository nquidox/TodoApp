package user

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"todoApp/api/service"
)

// CreateUserHandler     godoc
//
//	@Summary		Create user
//	@Description	Creates new user account.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			user	body		User	true	"Create new user"
//	@Success		200		{object}	User
//	@Failure		400		{object}	service.errorResponse	"Bad request"
//	@Failure		500		{object}	service.errorResponse	"Internal Server Error"
//	@Router			/user [post]
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	usr := User{}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		service.BadRequestResponse(w, service.BodyReadErr, err)
		log.Error(service.BodyReadErr, err)
		return
	}

	err = service.DeserializeJSON(data, &usr)
	if err != nil {
		service.UnprocessableEntity(w, service.JSONDeserializingErr, err)
		log.Error(service.JSONDeserializingErr, err)
		return
	}

	err = usr.CheckRequiredFields()
	if err != nil {
		service.BadRequestResponse(w, service.ValidationErr, err)
		log.Error(service.ValidationErr, err)
		return
	}

	err = usr.Create(Worker)
	if err != nil {
		service.InternalServerErrorResponse(w, service.UserCreateErr, err)
		log.Error(service.UserCreateErr, err)
		return
	}

	service.OkResponse(w, usr)

	log.WithFields(log.Fields{
		"id":       usr.UserUUID,
		"username": usr.Username,
		"email":    usr.Email,
	}).Info(service.UserCreateSuccess)
}

// ReadUserHandler     godoc
//
//	@Summary		Get user
//	@Description	Returns user info. UUID is optional for superusers.
//	@Tags			User
//	@Security		BasicAuth
//	@Produce		json
//	@Param			id	path		string	false	"uuid"
//	@Success		200	{object}	User
//	@Failure		400	{object}	service.errorResponse	"Bad request"
//	@Failure		401	{object}	service.errorResponse	"Unauthorized"
//	@Failure		500	{object}	service.errorResponse	"Internal Server Error"
//	@Router			/user/{id} [get]
func ReadUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	target := targetUUID(w, r)

	usr := User{UserUUID: target}
	err := usr.Read(Worker)
	if err != nil {
		log.Error(service.UserReadErr, err)
		service.InternalServerErrorResponse(w, service.UserReadErr, err)
		return
	}

	service.OkResponse(w, usr)
	log.WithFields(log.Fields{
		"id":       usr.UserUUID,
		"username": usr.Username,
		"email":    usr.Email,
	}).Info(service.UserReadSuccess)
}

// UpdateUserHandler     godoc
//
//	@Summary		Update user
//	@Description	Updates your account data. UUID is optional for superusers.
//	@Tags			User
//	@Security		BasicAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string		false	"uuid"
//	@Param			user	body		updateUser	true	"Update user"
//	@Success		200		{object}	service.errorResponse
//	@Failure		400		{object}	service.errorResponse	"Bad request"
//	@Failure		401		{object}	service.errorResponse	"Unauthorized"
//	@Failure		500		{object}	service.errorResponse	"Internal Server Error"
//	@Router			/user/{id} [put]
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	target := targetUUID(w, r)

	usr := updateUser{UserUUID: target}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		service.BadRequestResponse(w, service.BodyReadErr, err)
		log.Error(service.BodyReadErr, err)
		return
	}

	err = service.DeserializeJSON(data, &usr)
	if err != nil {
		service.UnprocessableEntity(w, service.JSONDeserializingErr, err)
		log.Error(service.JSONDeserializingErr, err)
		return
	}

	err = usr.Update(Worker)
	if err != nil {
		service.InternalServerErrorResponse(w, service.UserUpdateErr, err)
		log.Error(service.UserUpdateErr, err)
		return
	}

	service.OkResponse(w, service.DefaultResponse{
		ResultCode: 0,
		HttpCode:   http.StatusOK,
		Messages:   service.UserUpdateSuccess,
		Data:       nil,
	})

	log.WithFields(log.Fields{
		"id":       usr.UserUUID,
		"username": usr.Username,
		"email":    usr.Email,
	}).Info(service.UserUpdateSuccess)
}

// DeleteUserHandler     godoc
//
//	@Summary		Delete user
//	@Description	Deletes user account. UUID is optional for superusers.
//	@Tags			User
//	@Security		BasicAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	false	"uuid"
//	@Success		204
//	@Failure		400	{object}	service.errorResponse	"Bad request"
//	@Failure		401	{object}	service.errorResponse	"Unauthorized"
//	@Failure		500	{object}	service.errorResponse	"Internal Server Error"
//	@Router			/user/{id} [delete]
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	target := targetUUID(w, r)
	usr := User{UserUUID: target}

	err := usr.Delete(Worker)
	if err != nil {
		service.InternalServerErrorResponse(w, service.UserDeleteErr, err)
		return
	}

	service.OkResponse(w, service.DefaultResponse{
		ResultCode: 0,
		HttpCode:   http.StatusOK,
		Messages:   service.UserDeleteSuccess,
		Data:       nil,
	})

	log.WithFields(log.Fields{
		"id": usr.UserUUID,
	}).Info(service.UserDeleteSuccess)
}

func targetUUID(w http.ResponseWriter, r *http.Request) uuid.UUID {
	token, err := getTokenValue(r)
	if err != nil {
		log.Error(service.TokenReadErr, err.Error())
		service.BadRequestResponse(w, service.CookieReadErr, err)
		return uuid.Nil
	}

	t, userUUID := tokenIsValid(token)
	if !t {
		log.Error(service.TokenValidationErr, err)
		service.UnauthorizedResponse(w, "")
		return uuid.Nil
	}

	id := r.PathValue("id")
	parsedUUID := uuid.Nil

	if id != "" {
		parsedUUID, err = uuid.Parse(id)
		if err != nil {
			log.Warning(service.UUIDParseErr, err, "ignoring")
		}
	}

	usr := User{UserUUID: userUUID}
	err = usr.Read(Worker)
	if err != nil {
		log.Error(service.UserReadErr, err.Error())
		service.InternalServerErrorResponse(w, service.UserReadErr, err)
		return uuid.Nil
	}

	if usr.IsSuperuser && parsedUUID != uuid.Nil {
		return parsedUUID
	}

	return userUUID
}
