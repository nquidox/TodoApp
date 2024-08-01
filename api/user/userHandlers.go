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
//	@Failure		400		{object}	service.ErrorResponse	"Bad request"
//	@Failure		401		{object}	service.ErrorResponse	"Unauthorized"
//	@Failure		500		{object}	service.ErrorResponse	"Internal Server Error"
//	@Router			/user [post]
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	usr := User{}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusBadRequest,
			Messages:   "Error reading body: " + err.Error(),
			Data:       "",
		})
		return
	}

	err = service.DeserializeJSON(data, &usr)
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusInternalServerError,
			Messages:   "Error deserializing user: " + err.Error(),
			Data:       "",
		})
		return
	}

	err = usr.CheckRequiredFields()
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusBadRequest,
			Messages:   "Validation error: " + err.Error(),
			Data:       "",
		})
		return
	}

	err = usr.Create()
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusInternalServerError,
			Messages:   "Error creating user: " + err.Error(),
			Data:       "",
		})
		return
	}

	service.ServerResponse(w, usr)
}

// ReadUserHandler     godoc
//
//	@Summary		Get user
//	@Description	Get user info, uuid required
//	@Tags			User
//	@Security		CookieAuth
//	@Produce		json
//	@Param			id	path		string	true	"uuid"
//	@Success		200	{object}	User
//	@Failure		400	{object}	service.ErrorResponse	"Bad request"
//	@Failure		401	{object}	service.ErrorResponse	"Unauthorized"
//	@Failure		500	{object}	service.ErrorResponse	"Internal Server Error"
//	@Router			/user/{id} [get]
func ReadUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userUUID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusBadRequest,
			Messages:   "ID parse error: " + err.Error(),
			Data:       "",
		})
		return
	}

	err = Authorized(r, userUUID)
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusUnauthorized,
			Messages:   "Unauthorized",
			Data:       "",
		})
		return
	}

	usr := User{UserUUID: userUUID}
	err = usr.Read()
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusInternalServerError,
			Messages:   "User read error: " + err.Error(),
			Data:       "",
		})
		return
	}

	service.ServerResponse(w, usr)
}

// UpdateUserHandler     godoc
//
//	@Summary		Update user
//	@Description	Update user account
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			user	body		User	true	"Update user"
//	@Success		200		{object}	service.ErrorResponse
//	@Failure		400		{object}	service.ErrorResponse	"Bad request"
//	@Failure		401		{object}	service.ErrorResponse	"Unauthorized"
//	@Failure		500		{object}	service.ErrorResponse	"Internal Server Error"
//	@Router			/user/{id} [put]
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userUUID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusBadRequest,
			Messages:   "ID parse error: " + err.Error(),
			Data:       "",
		})
		return
	}

	usr := User{UserUUID: userUUID}

	err = Authorized(r, userUUID)
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusUnauthorized,
			Messages:   "Unauthorized",
			Data:       "",
		})
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusBadRequest,
			Messages:   "Error reading body: " + err.Error(),
			Data:       "",
		})
		return
	}

	err = service.DeserializeJSON(data, &usr)
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusInternalServerError,
			Messages:   "Error deserializing user: " + err.Error(),
			Data:       "",
		})
		return
	}

	err = usr.Update()
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusInternalServerError,
			Messages:   "User update error: " + err.Error(),
			Data:       "",
		})
		return
	}

	service.ServerResponse(w, service.ErrorResponse{
		ResultCode: 0,
		ErrorCode:  http.StatusOK,
		Messages:   "User updated successfully",
		Data:       "",
	})
}

// DeleteUserHandler     godoc
//
//	@Summary		Delete user
//	@Description	Delete user account
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			user	body		User	true	"Delete user"
//	@Success		200		{object}	service.ErrorResponse
//	@Failure		400		{object}	service.ErrorResponse	"Bad request"
//	@Failure		401		{object}	service.ErrorResponse	"Unauthorized"
//	@Failure		500		{object}	service.ErrorResponse	"Internal Server Error"
//	@Router			/user/{id} [delete]
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userUUID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusBadRequest,
			Messages:   "ID parse error: " + err.Error(),
			Data:       "",
		})
		return
	}

	usr := User{UserUUID: userUUID}

	err = Authorized(r, userUUID)
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusUnauthorized,
			Messages:   "Unauthorized",
			Data:       "",
		})
		return
	}

	err = usr.Delete()
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusInternalServerError,
			Messages:   "User delete error: " + err.Error(),
			Data:       "",
		})
		return
	}

	service.ServerResponse(w, service.ErrorResponse{
		ResultCode: 0,
		ErrorCode:  http.StatusOK,
		Messages:   "User deleted successfully",
		Data:       "",
	})
}
