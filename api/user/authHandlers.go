package user

import (
	"github.com/google/uuid"
	"io"
	"net/http"
	"todoApp/api/service"
)

func MeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token, err := r.Cookie("token")
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusUnauthorized,
			Messages:   "Unauthorized",
			Data:       "",
		})
		return
	}

	tokenString := token.Value
	s := Session{}
	err = DB.Where("token = ?", tokenString).First(&s).Error
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusUnauthorized,
			Messages:   "Unauthorized",
			Data:       "",
		})
		return
	}

	usr := User{Uuid: s.Uuid}
	err = DB.Model(&usr).Where("uuid = ?", s.Uuid).First(&usr).Error
	if err != nil {
		service.ServerResponse(w, service.ErrorResponse{
			ResultCode: 1,
			ErrorCode:  http.StatusInternalServerError,
			Messages:   "Error reading user from DB: " + err.Error(),
			Data:       "",
		})
		return
	}

	type shortUser struct {
		Uuid     uuid.UUID `json:"id"`
		Email    string    `json:"email"`
		Username string    `json:"login"`
	}

	service.ServerResponse(w, service.ErrorResponse{
		ResultCode: 0,
		ErrorCode:  200,
		Messages:   "",
		Data: shortUser{
			Uuid:     usr.Uuid,
			Email:    usr.Email,
			Username: usr.Username,
		},
	})
}

// LoginHandler     godoc
//
//	@Summary		Log in
//	@Description	Success login gives you a cookie with access token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		User	true	"Login"
//	@Success		200		{object}	service.ErrorResponse
//	@Failure		400		{object}	service.ErrorResponse	"Bad request"
//	@Failure		401		{object}	service.ErrorResponse	"Unauthorized"
//	@Failure		500		{object}	service.ErrorResponse	"Internal Server Error"
//	@Router			/login [post]
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

// LogoutHandler     godoc
//
//	@Summary		Log out
//	@Description	Log out and invalidate access token
//	@Tags			Auth
//	@Success		200
//	@Failure		400	{object}	service.ErrorResponse	"Bad request"
//	@Failure		401	{object}	service.ErrorResponse	"Unauthorized"
//	@Failure		500	{object}	service.ErrorResponse	"Internal Server Error"
//	@Router			/logout [get]
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
