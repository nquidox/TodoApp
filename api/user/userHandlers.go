package user

import (
	"github.com/google/uuid"
	"io"
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

	userUUID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		service.BadRequestResponse(w, service.ParseErr, err)
		return
	}

	err = Authorized(r, userUUID)
	if err != nil {
		service.UnauthorizedResponse(w, "")
		return
	}

	usr := User{UserUUID: userUUID}
	err = usr.Read()
	if err != nil {
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
//	@Param			id		path		string	false	"uuid"
//	@Param			user	body		User	true	"Update user"
//	@Success		200		{object}	service.errorResponse
//	@Failure		400		{object}	service.errorResponse	"Bad request"
//	@Failure		401		{object}	service.errorResponse	"Unauthorized"
//	@Failure		500		{object}	service.errorResponse	"Internal Server Error"
//	@Router			/user/{id} [put]
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userUUID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		service.BadRequestResponse(w, service.ParseErr, err)
		return
	}

	if userUUID == uuid.Nil {
		token, err := r.Cookie("token")
		if err != nil {
			service.UnauthorizedResponse(w, "")
			return
		}

		userUUID, err = getUserUUIDFromToken(token.Value)
		if err != nil {
			service.UnauthorizedResponse(w, "")
			return
		}
	}

	usr := User{UserUUID: userUUID}

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
//	@Param			id		path		string	false	"uuid"
//	@Param			user	body		User	true	"Delete user"
//	@Success		200		{object}	service.errorResponse
//	@Failure		400		{object}	service.errorResponse	"Bad request"
//	@Failure		401		{object}	service.errorResponse	"Unauthorized"
//	@Failure		500		{object}	service.errorResponse	"Internal Server Error"
//	@Router			/user/{id} [delete]
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userUUID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		service.BadRequestResponse(w, service.ParseErr, err)
		return
	}

	usr := User{UserUUID: userUUID}

	err = Authorized(r, userUUID)
	if err != nil {
		service.UnauthorizedResponse(w, "")
		return
	}

	err = usr.Delete()
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
