package user

import (
	"github.com/google/uuid"
	"io"
	"net/http"
	"todoApp/api/service"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data, err := io.ReadAll(r.Body)
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusBadRequest,
			Messages:   "JSON read error: " + err.Error(),
			Data:       "",
		})
	}

	usr := User{}
	err = service.DeserializeJSON(data, &usr)
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusInternalServerError,
			Messages:   "JSON parsing error: " + err.Error(),
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

	getUsr, err := getUser(usr.Email)
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusBadRequest,
			Messages:   "Email is incorrect",
			Data:       "",
		})
		return
	}

	err = comparePasswords(getUsr.Password, usr.Password)
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusBadRequest,
			Messages:   "Password is incorrect",
		})
		return
	}

	cookie, err := createSession(getUsr.Uuid)
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusInternalServerError,
			Messages:   err.Error(),
			Data:       "",
		})
		return
	}

	type uuidOnly struct {
		UUID uuid.UUID `json:"userId"`
	}

	http.SetCookie(w, &cookie)
	service.ServerResponse(w, service.ErrorResponse{
		ResultCode: 0,
		ErrorCode:  http.StatusOK,
		Messages:   "",
		Data:       uuidOnly{UUID: getUsr.Uuid},
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusBadRequest,
			Messages:   err.Error(),
			Data:       "",
		})
		return
	}

	err = dropSession(cookie.Value)
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusBadRequest,
			Messages:   err.Error(),
			Data:       "",
		})
		return
	}
}

func getUser(email string) (User, error) {
	var user User
	err := DB.Where("email = ?", email).First(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}
