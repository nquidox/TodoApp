package user

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"todoApp/api/service"
)

// meFunc     godoc
//
//	@Summary		meFunc request
//	@Description	meFunc request
//	@Security		BasicAuth
//	@Tags			Auth
//	@Produce		json
//	@Success		200	{object}	meResponse
//	@Failure		400	{object}	service.errorResponse	"Bad request"
//	@Failure		401	{object}	service.errorResponse	"Unauthorized"
//	@Failure		500	{object}	service.errorResponse	"Internal Server Error"
//	@Router			/meFunc [get]
func meFunc(s *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		token, err := r.Cookie(service.SessionTokenName)
		if err != nil {
			service.UnauthorizedResponse(w, "")
			log.Error(service.TokenReadErr, err)
			return
		}

		session := Session{Token: token.Value}
		err = session.Read(s.DbWorker)
		if err != nil {
			service.UnauthorizedResponse(w, "")
			log.Error(service.TokenValidationErr, err)
			return
		}

		me := meModel{UserUUID: session.UserUuid}
		err = me.Read(s.DbWorker)
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
		}).Info("/meFunc ", service.UserReadSuccess)
	}
}

// loginFunc     godoc
//
//	@Summary		Log in
//	@Description	Success loginFunc gives you a cookie with access token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		loginUserModel	true	"loginFunc"
//	@Success		200		{object}	service.errorResponse
//	@Failure		400		{object}	service.errorResponse	"Bad request"
//	@Failure		401		{object}	service.errorResponse	"Unauthorized"
//	@Failure		500		{object}	service.errorResponse	"Internal Server Error"
//	@Router			/loginFunc [post]
func loginFunc(s *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		err = getUsr.Read(s.DbWorker)
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

		var session Session
		cookie, err := session.Create(s.DbWorker, getUsr.UserUUID, s.Salt)
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
}

// logoutFunc     godoc
//
//	@Summary		Log out
//	@Description	Log out and invalidate access token
//	@Tags			Auth
//	@Security		ApiKeyAuth
//	@Success		200
//	@Failure		400	{object}	service.errorResponse	"Bad request"
//	@Failure		401	{object}	service.errorResponse	"Unauthorized"
//	@Failure		500	{object}	service.errorResponse	"Internal Server Error"
//	@Router			/logoutFunc [get]
func logoutFunc(s *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(service.SessionTokenName)
		if err != nil {
			service.BadRequestResponse(w, service.CookieReadErr, err)
			log.Error(service.CookieReadErr, err)
			return
		}

		session := Session{Token: cookie.Value}
		err = session.Delete(s.DbWorker)
		if err != nil {
			service.UnauthorizedResponse(w, service.InvalidTokenErr)
			log.Error(service.InvalidTokenErr)
			return
		}

		log.WithFields(log.Fields{
			"session": cookie.Value,
		}).Info(service.LogoutSuccess)
	}
}
