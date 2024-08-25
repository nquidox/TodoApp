package user

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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

	token, err := r.Cookie(service.SessionTokenName)
	if err != nil {
		service.UnauthorizedResponse(w, "")
		log.Error(service.TokenReadErr, err)
		return
	}

	s := Session{Token: token.Value}
	err = s.Read(dbw)
	if err != nil {
		service.UnauthorizedResponse(w, "")
		log.Error(service.TokenValidationErr, err)
		return
	}

	me := meModel{UserUUID: s.UserUuid}
	err = me.Read(dbw)
	if err != nil {
		service.BadRequestResponse(w, service.UserReadErr, err)
		log.Error(service.UserReadErr, err)
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

	log.WithFields(log.Fields{
		"id":       me.UserUUID,
		"username": me.Username,
	}).Info("/me ", service.UserReadSuccess)
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
		log.Error(service.JSONReadErr, err)
		return
	}

	usr := loginUserModel{}
	err = service.DeserializeJSON(data, &usr)
	if err != nil {
		service.UnprocessableEntityResponse(w, service.JSONDeserializingErr, err)
		log.Error(service.JSONDeserializingErr, err)
		return
	}

	err = usr.CheckRequiredFields()
	if err != nil {
		service.BadRequestResponse(w, service.ValidationErr, err)
		log.Error(service.ValidationErr, err)
		return
	}

	getUsr := User{Email: usr.Email}
	err = getUsr.Read(dbw)
	if err != nil {
		service.BadRequestResponse(w, service.EmailErr, err)
		log.Error(service.EmailErr, err)
		return
	}

	err = comparePasswords(getUsr.Password, usr.Password)
	if err != nil {
		service.BadRequestResponse(w, service.PasswordErr, err)
		log.Error(service.PasswordErr, err)
		return
	}

	var s Session
	cookie, err := s.Create(dbw, getUsr.UserUUID)
	if err != nil {
		service.InternalServerErrorResponse(w, service.SessionCreateErr, err)
		log.Error(service.SessionCreateErr, err)
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

	log.WithFields(log.Fields{
		"id":       getUsr.UserUUID,
		"username": getUsr.Username,
	}).Info(service.LoginSuccess)
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
	cookie, err := r.Cookie(service.SessionTokenName)
	if err != nil {
		service.BadRequestResponse(w, service.CookieReadErr, err)
		log.Error(service.CookieReadErr, err)
		return
	}

	s := Session{Token: cookie.Value}
	err = s.Delete(dbw)
	if err != nil {
		service.UnauthorizedResponse(w, service.InvalidTokenErr)
		log.Error(service.InvalidTokenErr)
		return
	}

	log.WithFields(log.Fields{
		"session": cookie.Value,
	}).Info(service.LogoutSuccess)
}
