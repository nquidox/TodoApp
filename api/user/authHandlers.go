package user

import (
	"github.com/google/uuid"
	"io"
	"net/http"
	"todoApp/api/service"
)

// MeHandler     godoc
//
//	@Summary		Me request
//	@Description	Me request
//	@Security		BasicAuth
//	@Tags			Auth
//	@Produce		json
//	@Success		200	{object}	meResponse
//	@Failure		400	{object}	service.errorResponse	"Bad request"
//	@Failure		401	{object}	service.errorResponse	"Unauthorized"
//	@Failure		500	{object}	service.errorResponse	"Internal Server Error"
//	@Router			/me [get]
func MeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token, err := r.Cookie("token")
	if err != nil {
		service.UnauthorizedResponse(w, "")
		return
	}

	s := Session{}
	err = DB.Where("token = ?", token.Value).First(&s).Error
	if err != nil {
		service.UnauthorizedResponse(w, "")
		return
	}

	me := meModel{UserUUID: s.UserUuid}
	err = DB.Model(User{}).Where("user_uuid = ?", s.UserUuid).First(&me).Error
	if err != nil {
		service.InternalServerErrorResponse(w, service.DBReadErr, err)
		return
	}

	service.OkResponse(w, meResponse{
		ResultCode: 0,
		HttpCode:   200,
		Messages:   nil,
		Data: meModel{
			UserUUID: me.UserUUID,
			Email:    me.Email,
			Username: me.Username,
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
//	@Param			user	body		loginUserModel	true	"Login"
//	@Success		200		{object}	service.errorResponse
//	@Failure		400		{object}	service.errorResponse	"Bad request"
//	@Failure		401		{object}	service.errorResponse	"Unauthorized"
//	@Failure		500		{object}	service.errorResponse	"Internal Server Error"
//	@Router			/login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data, err := io.ReadAll(r.Body)
	if err != nil {
		service.BadRequestResponse(w, service.JSONReadErr, err)
		return
	}

	usr := loginUserModel{}
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

	getUsr, err := getUser(usr.Email)
	if err != nil {
		service.BadRequestResponse(w, service.EmailErr, err)
		return
	}

	err = comparePasswords(getUsr.Password, usr.Password)
	if err != nil {
		service.BadRequestResponse(w, service.PasswordErr, err)
		return
	}

	cookie, err := createSession(getUsr.UserUUID)
	if err != nil {
		service.InternalServerErrorResponse(w, service.SessionCreateErr, err)
		return
	}

	type uuidOnly struct {
		UUID uuid.UUID `json:"userId"`
	}

	http.SetCookie(w, &cookie)
	service.OkResponse(w, service.DefaultResponse{
		ResultCode: 0,
		HttpCode:   http.StatusOK,
		Messages:   "",
		Data:       uuidOnly{UUID: getUsr.UserUUID},
	})
}

// LogoutHandler     godoc
//
//	@Summary		Log out
//	@Description	Log out and invalidate access token
//	@Tags			Auth
//	@Security		ApiKeyAuth
//	@Success		200
//	@Failure		400	{object}	service.errorResponse	"Bad request"
//	@Failure		401	{object}	service.errorResponse	"Unauthorized"
//	@Failure		500	{object}	service.errorResponse	"Internal Server Error"
//	@Router			/logout [get]
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		service.BadRequestResponse(w, service.CookieReadErr, err)
		return
	}

	err = dropSession(cookie.Value)
	if err != nil {
		service.InternalServerErrorResponse(w, service.SessionCloseErr, err)
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
