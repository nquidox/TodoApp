package user

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
	"time"
	"todoApp/api/service"
)

// createUserFunc     godoc
//
//	@Summary		Create user
//	@Description	Creates new user account.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			user	body		User					true	"Create new user"
//	@Success		200		{object}	User					"OK"
//	@Failure		400		{object}	service.errorResponse	"Bad request"
//	@Failure		409		{object}	service.errorResponse	"Conflict"
//	@Failure		422		{object}	service.errorResponse	"Unprocessable entity"
//	@Failure		500		{object}	service.errorResponse	"Internal server error"
//	@Router			/user [post]
func createUserFunc(s *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usr := User{}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(service.BodyReadErr, err)
			service.BadRequestResponse(w, service.BodyReadErr, err)
			return
		}

		err = service.DeserializeJSON(data, &usr)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			log.Error(service.JSONDeserializingErr, err)
			service.UnprocessableEntityResponse(w, service.JSONDeserializingErr, err)
			return
		}

		usr.Email = strings.ToLower(usr.Email)

		err = usr.CheckRequiredFields()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(service.ValidationErr, err)
			service.BadRequestResponse(w, service.ValidationErr, err)
			return
		}

		err = usr.emailExists(s.DbWorker)
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			log.Error(service.ConflictErr, err)
			service.ConflictResponse(w, service.ConflictErr)
			return
		}

		verificationKey, err := generateEmailVerificationKey()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(service.VerificationKeyErr, err)
			service.InternalServerErrorResponse(w, service.VerificationKeyErr, err)
			return
		}

		usr.EmailVerificationKey = verificationKey
		usr.EmailKeyCreatedAt = time.Now()

		err = usr.Create(s.DbWorker)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(service.UserCreateErr, err)
			service.InternalServerErrorResponse(w, service.UserCreateErr, err)
			return
		}

		err = sendVerificationEmail(usr.Email, verificationKey, s)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(service.EmailSendErr, err)
			service.InternalServerErrorResponse(w, service.EmailSendErr, err)
			return
		}
		log.Debug("Verification link sent on create user")

		w.WriteHeader(http.StatusOK)
		log.WithFields(log.Fields{
			"id":       usr.UserUUID,
			"username": usr.Username,
			"email":    usr.Email,
		}).Info(service.UserCreateSuccess)
		service.OkResponse(w, usr)
	}
}

// readUserFunc     godoc
//
//	@Security		BasicAuth
//	@Summary		Get user
//	@Description	Returns user info. UUID is optional for superusers.
//	@Tags			User
//	@Produce		json
//	@Param			id	path		string					false	"uuid"
//	@Success		200	{object}	User					"OK"
//	@Failure		400	{object}	service.errorResponse	"Bad request"
//	@Failure		401	{object}	service.errorResponse	"Unauthorized"
//	@Failure		404	{object}	service.errorResponse	"Not Found"
//	@Failure		500	{object}	service.errorResponse	"Internal Server Error"
//	@Router			/user/{id} [get]
func readUserFunc(s *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		target := targetUUID(w, r, s)

		usr := readUser{UserUUID: target}
		err := usr.Read(s.DbWorker)
		if err != nil {
			if err.Error() == "404" {
				w.WriteHeader(http.StatusNotFound)
				log.Error(service.DBNotFound)
				service.NotFoundResponse(w, "")
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			service.InternalServerErrorResponse(w, service.UserReadErr, err)
			log.Error(service.UserReadErr, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.WithFields(log.Fields{
			"id":       usr.UserUUID,
			"username": usr.Username,
			"email":    usr.Email,
		}).Info(service.UserReadSuccess)
		service.OkResponse(w, usr)
	}
}

// updateUserFunc     godoc
//
//	@Security		BasicAuth
//	@Summary		Update user
//	@Description	Updates your account data. UUID is optional for superusers.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					false	"uuid"
//	@Param			user	body		updateUser				true	"Update user"
//	@Success		200		{object}	service.DefaultResponse	"OK"
//	@Failure		400		{object}	service.errorResponse	"Bad request"
//	@Failure		401		{object}	service.errorResponse	"Unauthorized"
//	@Failure		404		{object}	service.errorResponse	"Not Found"
//	@Failure		422		{object}	service.errorResponse	"Unprocessable entity"
//	@Failure		500		{object}	service.errorResponse	"Internal server error"
//	@Router			/user/{id} [put]
func updateUserFunc(s *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		target := targetUUID(w, r, s)
		usr := updateUser{UserUUID: target}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(service.BodyReadErr, err)
			service.BadRequestResponse(w, service.BodyReadErr, err)
			return
		}

		err = service.DeserializeJSON(data, &usr)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			log.Error(service.JSONDeserializingErr, err)
			service.UnprocessableEntityResponse(w, service.JSONDeserializingErr, err)
			return
		}

		usr.Email = strings.ToLower(usr.Email)

		err = usr.Update(s.DbWorker)
		if err != nil {
			if err.Error() == "404" {
				w.WriteHeader(http.StatusNotFound)
				log.WithFields(log.Fields{
					"id":       usr.UserUUID,
					"username": usr.Username,
					"email":    usr.Email,
				}).Error(service.DBNotFound)
				service.NotFoundResponse(w, "")
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(service.UserUpdateErr, err)
			service.InternalServerErrorResponse(w, service.UserUpdateErr, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.WithFields(log.Fields{
			"id":       usr.UserUUID,
			"username": usr.Username,
			"email":    usr.Email,
		}).Info(service.UserUpdateSuccess)

		service.OkResponse(w, service.DefaultResponse{
			ResultCode: 0,
			HttpCode:   http.StatusOK,
			Messages:   service.UserUpdateSuccess,
			Data:       nil,
		})
	}
}

// deleteUserFunc     godoc
//
//	@Security		BasicAuth
//	@Summary		Delete user
//	@Description	Deletes user account. UUID is optional for superusers.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string					false	"uuid"
//	@Success		200	{object}	service.DefaultResponse	"OK"
//	@Failure		400	{object}	service.errorResponse	"Bad request"
//	@Failure		401	{object}	service.errorResponse	"Unauthorized"
//	@Failure		404	{object}	service.errorResponse	"Not Found"
//	@Failure		500	{object}	service.errorResponse	"Internal Server Error"
//	@Router			/user/{id} [delete]
func deleteUserFunc(s *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		target := targetUUID(w, r, s)
		usr := User{UserUUID: target}

		err := usr.Delete(s.DbWorker)
		if err != nil {
			if err.Error() == "404" {
				w.WriteHeader(http.StatusNotFound)
				service.NotFoundResponse(w, "")
				log.Error(service.DBNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(service.UserDeleteErr, err)
			service.InternalServerErrorResponse(w, service.UserDeleteErr, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.WithFields(log.Fields{
			"id": usr.UserUUID,
		}).Info(service.UserDeleteSuccess)

		service.OkResponse(w, service.DefaultResponse{
			ResultCode: 0,
			HttpCode:   http.StatusNoContent,
			Messages:   service.UserDeleteSuccess,
			Data:       nil,
		})
	}
}

func targetUUID(w http.ResponseWriter, r *http.Request, s *Service) uuid.UUID {
	token, err := r.Cookie(service.SessionTokenName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(service.TokenReadErr, err.Error())
		service.BadRequestResponse(w, service.CookieReadErr, err)
		return uuid.Nil
	}

	session := Session{Token: token.Value}
	err = session.Read(s.DbWorker)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Error(service.TokenValidationErr, err)
		service.UnauthorizedResponse(w, "")
		return uuid.Nil
	}

	id := r.PathValue("id")
	parsedUUID := uuid.Nil

	if id != "" {
		parsedUUID, err = uuid.Parse(id)
		if err != nil {
			log.Warning(service.UUIDParseErr, err, ", ignoring")
		}
	}

	usr := User{UserUUID: session.UserUuid}
	err = usr.Read(s.DbWorker)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(service.UserReadErr, err.Error())
		service.InternalServerErrorResponse(w, service.UserReadErr, err)
		return uuid.Nil
	}

	if usr.IsSuperuser && parsedUUID != uuid.Nil {
		return parsedUUID
	}

	return session.UserUuid
}
