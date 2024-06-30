package user

import (
	"github.com/google/uuid"
	"io"
	"net/http"
	"todoApp/api/service"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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

	usr := User{Uuid: userUUID}
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

	usr := User{Uuid: userUUID}

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

	usr := User{Uuid: userUUID}

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
