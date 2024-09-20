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
			log.Error(service.TokenReadErr, err)
			service.UnauthorizedResponse(w, "")
			return
		}

		session := Session{Token: token.Value}
		err = session.Read(s.DbWorker)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Error(service.TokenValidationErr, err)
			service.UnauthorizedResponse(w, "")
			return
		}

		me := meModel{UserUUID: session.UserUuid}
		err = me.Read(s.DbWorker)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(service.UserReadErr, err)
			service.BadRequestResponse(w, service.UserReadErr, err)
			return
		}

		w.WriteHeader(http.StatusOK)

		log.WithFields(log.Fields{
			"id":       me.UserUUID,
			"username": me.Username,
		}).Info("/me ", service.UserReadSuccess)

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
			log.Error(service.JSONReadErr, err)
			service.BadRequestResponse(w, service.JSONReadErr, err)
			return
		}

		usr := loginUserModel{}
		err = service.DeserializeJSON(data, &usr)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			log.Error(service.JSONDeserializingErr, err)
			service.UnprocessableEntityResponse(w, service.JSONDeserializingErr, err)
			return
		}

		err = usr.CheckRequiredFields()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(service.ValidationErr, err)
			service.BadRequestResponse(w, service.ValidationErr, err)
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
			log.Error(service.EmailErr, err)
			service.BadRequestResponse(w, service.EmailErr, err)
			return
		}

		if !getUsr.EmailVerified {
			w.WriteHeader(http.StatusForbidden)
			log.Error(service.EmailNotVerified, err)
			service.ForbiddenResponse(w, service.EmailNotVerified)
			return
		}

		err = comparePasswords(getUsr.Password, usr.Password)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(service.PasswordErr, err)
			service.BadRequestResponse(w, service.PasswordErr, err)
			return
		}

		var session Session
		cookie, err := session.Create(s.DbWorker, getUsr.UserUUID, s.Salt, r.UserAgent())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(service.SessionCreateErr, err)
			service.InternalServerErrorResponse(w, service.SessionCreateErr, err)
			return
		}

		type uuidOnly struct {
			UUID uuid.UUID `json:"userId"`
		}

		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)

		log.WithFields(log.Fields{
			"id":       getUsr.UserUUID,
			"username": getUsr.Username,
		}).Info(service.LoginSuccess)

		service.OkResponse(w, service.DefaultResponse{
			ResultCode: 0,
			HttpCode:   http.StatusOK,
			Messages:   "",
			Data:       uuidOnly{UUID: getUsr.UserUUID},
		})
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
//	@Router			/logout [post]
func logoutFunc(s *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(service.SessionTokenName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(service.CookieReadErr, err)
			service.BadRequestResponse(w, service.CookieReadErr, err)
			return
		}

		session := Session{Token: cookie.Value}
		err = session.Delete(s.DbWorker)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Error(service.InvalidTokenErr)
			service.UnauthorizedResponse(w, service.InvalidTokenErr)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		log.WithFields(log.Fields{
			"session": cookie.Value,
		}).Info(service.LogoutSuccess)
	}
}

// getAllSessionsFunc godoc
//
//	@Security		BasicAuth
//	@Summary		Get all user's sessions
//	@Description	Requests all user's sessions
//	@Tags			Session
//	@Produce		json
//	@Success		200	{array}		Session					"OK"
//	@Success		204	{array}		service.DefaultResponse	"No Content"
//	@Failure		401	{object}	service.errorResponse	"Unauthorized"
//	@Failure		500	{object}	service.errorResponse	"Internal server error"
//	@Router			/getAllSessions [get]
func getAllSessionsFunc(s *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		token, err := r.Cookie(service.SessionTokenName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Error(service.TokenReadErr, err)
			service.UnauthorizedResponse(w, "")
			return
		}

		session := Session{Token: token.Value}
		err = session.Read(s.DbWorker)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Error(service.TokenValidationErr, err)
			service.UnauthorizedResponse(w, "")
			return
		}

		list, err := session.ReadAll(s.DbWorker)
		if err != nil {
			if err.Error() == "404" {
				w.WriteHeader(http.StatusNoContent)
				log.Info(service.NoContent)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(service.DBReadErr, err)
			service.InternalServerErrorResponse(w, service.DBReadErr, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info(service.SessionsReadSuccess)
		service.OkResponse(w, list)
	}
}

// closeOtherSessionsFunc     godoc
//
//	@Security		BasicAuth
//	@Summary		Close all other sessions
//	@Description	Closes all other active sessions except current. To close this session use /logout endpoint in auth block.
//	@Tags			Session
//	@Success		204
//	@Failure		400	{object}	service.errorResponse	"Bad request"
//	@Failure		401	{object}	service.errorResponse	"Unauthorized"
//	@Failure		500	{object}	service.errorResponse	"Internal server error"
//	@Router			/closeOtherSessions [post]
func closeOtherSessionsFunc(s *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		token, err := r.Cookie(service.SessionTokenName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(service.CookieReadErr, err)
			service.BadRequestResponse(w, service.CookieReadErr, err)
			return
		}

		session := Session{Token: token.Value}
		err = session.Read(s.DbWorker)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Error(service.TokenValidationErr, err)
			service.UnauthorizedResponse(w, "")
			return
		}

		err = session.DeleteAllExceptOne(s.DbWorker, session.Token)
		if err != nil {
			if err.Error() == "404" {
				w.WriteHeader(http.StatusNoContent)
				log.Info(service.NoContent)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(service.DBReadErr, err)
			service.InternalServerErrorResponse(w, service.DBReadErr, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		log.Info(service.SessionsCloseSuccess)
	}
}
