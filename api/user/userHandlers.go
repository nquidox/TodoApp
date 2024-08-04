package user

import (
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"todoApp/api/service"
)

// CreateUserHandler     godoc
//
//	@Summary		Create user
//	@Description	Create new user account
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
		return
	}

	err = service.DeserializeJSON(data, &usr)
	if err != nil {
		service.UnprocessableEntity(w, service.JSONDeserializingErr, err)
		return
	}

	err = usr.CheckRequiredFields()
	if err != nil {
		service.BadRequestResponse(w, service.ValidationErr, err)
		return
	}

	err = usr.Create()
	if err != nil {
		service.InternalServerErrorResponse(w, service.UserCreateErr, err)
		return
	}

	service.OkResponse(w, usr)
}

// ReadUserHandler     godoc
//
//	@Summary		Get user
//	@Description	Get user info, uuid required
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
	err := usr.Read()
	if err != nil {
		log.Println("Error parsed reading user", err.Error())
		service.InternalServerErrorResponse(w, service.UserReadErr, err)
		return
	}

	service.OkResponse(w, usr)
}

// UpdateUserHandler     godoc
//
//	@Summary		Update user
//	@Description	Update your account data
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
		return
	}

	err = service.DeserializeJSON(data, &usr)
	if err != nil {
		service.UnprocessableEntity(w, service.JSONDeserializingErr, err)
		return
	}

	err = usr.Update()
	if err != nil {
		service.InternalServerErrorResponse(w, service.UserUpdateErr, err)
		return
	}

	service.OkResponse(w, service.DefaultResponse{
		ResultCode: 0,
		HttpCode:   http.StatusOK,
		Messages:   service.UpdateOk,
		Data:       nil,
	})
}

// DeleteUserHandler     godoc
//
//	@Summary		Delete user
//	@Description	Delete user account
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

	err := usr.Delete()
	if err != nil {
		service.InternalServerErrorResponse(w, service.UserDeleteErr, err)
		return
	}

	service.OkResponse(w, service.DefaultResponse{
		ResultCode: 0,
		HttpCode:   http.StatusOK,
		Messages:   service.DeleteOk,
		Data:       nil,
	})
}

func targetUUID(w http.ResponseWriter, r *http.Request) uuid.UUID {
	token, err := getTokenValue(r)
	if err != nil {
		log.Println("Error getting token", err.Error())
		service.BadRequestResponse(w, service.CookieReadErr, err)
		return uuid.Nil
	}

	t, userUUID := tokenIsValid(token)
	if !t {
		log.Println("Error validating token", err)
		service.UnauthorizedResponse(w, "")
		return uuid.Nil
	}

	id := r.PathValue("id")
	parsedUUID := uuid.Nil

	if id != "" {
		parsedUUID, err = uuid.Parse(id)
		if err != nil {
			log.Println("Error parsing uuid", err, "Ignoring")
		}
	}

	usr := User{UserUUID: userUUID}
	err = usr.Read()
	if err != nil {
		log.Println("Error reading user", err.Error())
		service.InternalServerErrorResponse(w, service.UserReadErr, err)
		return uuid.Nil
	}

	if usr.IsSuperuser && parsedUUID != uuid.Nil {
		return parsedUUID
	}

	return userUUID
}
