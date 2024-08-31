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
//	@Security		BasicAuth
//	@Summary		me request
//	@Description	me request
//	@Tags			Auth
//	@Produce		json
//	@Success		200	{object}	meResponse				"OK"
//	@Failure		400	{object}	service.errorResponse	"Bad request"
//	@Failure		401	{object}	service.errorResponse	"Unauthorized"
//	@Failure		500	{object}	service.errorResponse	"Internal server error"
//	@Router			/me [get]
func meFunc(s *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		token, err := r.Cookie(service.SessionTokenName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			service.UnauthorizedResponse(w, "")
			log.Error(service.TokenReadErr, err)
			return
		}

		session := Session{Token: token.Value}
		err = session.Read(s.DbWorker)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			service.UnauthorizedResponse(w, "")
			log.Error(service.TokenValidationErr, err)
			return
		}

		me := meModel{UserUUID: session.UserUuid}
		err = me.Read(s.DbWorker)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			service.BadRequestResponse(w, service.UserReadErr, err)
			log.Error(service.UserReadErr, err)
			return
		}

		w.WriteHeader(http.StatusOK)
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
}

// loginFunc     godoc
//
//	@Summary		Log in
//	@Description	Success login gives you a cookie with access token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		loginUserModel			true	"login"
//	@Success		200		{object}	service.DefaultResponse	"OK"
//	@Failure		400		{object}	service.errorResponse	"Bad request"
//	@Failure		404		{object}	service.errorResponse	"Not found"
//	@Failure		500		{object}	service.errorResponse	"Internal server error"
//	@Router			/login [post]
func loginFunc(s *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		data, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			service.BadRequestResponse(w, service.JSONReadErr, err)
			log.Error(service.JSONReadErr, err)
			return
		}

		usr := loginUserModel{}
		err = service.DeserializeJSON(data, &usr)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			service.UnprocessableEntityResponse(w, service.JSONDeserializingErr, err)
			log.Error(service.JSONDeserializingErr, err)
			return
		}

		err = usr.CheckRequiredFields()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			service.BadRequestResponse(w, service.ValidationErr, err)
			log.Error(service.ValidationErr, err)
			return
		}

		getUsr := User{Email: usr.Email}
		err = getUsr.Read(s.DbWorker)
		if err != nil {
			if err.Error() == "404" {
				w.WriteHeader(http.StatusNotFound)
				log.Debug(service.EmailNotFoundErr)
				service.NotFoundResponse(w, service.EmailNotFoundErr)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			service.BadRequestResponse(w, service.EmailErr, err)
			log.Error(service.EmailErr, err)
			return
		}

		if !getUsr.EmailVerified {
			w.WriteHeader(http.StatusForbidden)
			service.ForbiddenResponse(w, service.EmailNotVerified)
			log.Error(service.EmailNotVerified, err)
			return
		}

		err = comparePasswords(getUsr.Password, usr.Password)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			service.BadRequestResponse(w, service.PasswordErr, err)
			log.Error(service.PasswordErr, err)
			return
		}

		var session Session
		cookie, err := session.Create(s.DbWorker, getUsr.UserUUID, s.Salt)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			service.InternalServerErrorResponse(w, service.SessionCreateErr, err)
			log.Error(service.SessionCreateErr, err)
			return
		}

		type uuidOnly struct {
			UUID uuid.UUID `json:"userId"`
		}

		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
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
//	@Security		BasicAuth
//	@Summary		Log out
//	@Description	Log out and invalidate access token
//	@Tags			Auth
//	@Success		204
//	@Failure		400	{object}	service.errorResponse	"Bad request"
//	@Failure		401	{object}	service.errorResponse	"Unauthorized"
//	@Failure		500	{object}	service.errorResponse	"Internal server error"
//	@Router			/logout [get]
func logoutFunc(s *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(service.SessionTokenName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			service.BadRequestResponse(w, service.CookieReadErr, err)
			log.Error(service.CookieReadErr, err)
			return
		}

		session := Session{Token: cookie.Value}
		err = session.Delete(s.DbWorker)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			service.UnauthorizedResponse(w, service.InvalidTokenErr)
			log.Error(service.InvalidTokenErr)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		log.WithFields(log.Fields{
			"session": cookie.Value,
		}).Info(service.LogoutSuccess)
	}
}
