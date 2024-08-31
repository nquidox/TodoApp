package user

import (
	"crypto/tls"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"net/http"
	"time"
	"todoApp/api/service"
)

// emailFunc     godoc
//
//	@Summary		Verify email
//	@Description	Verify email by key
//	@Tags			Email
//	@Produce		json
//	@Param			key	path		string					true	"key"
//	@Success		204	{object}	service.DefaultResponse	"No Content"
//	@Failure		400	{object}	service.errorResponse	"Bad request"
//	@Failure		404	{object}	service.errorResponse	"Not Found"
//	@Failure		500	{object}	service.errorResponse	"Internal Server Error"
//	@Router			/verifyEmail/{key} [post]
func emailFunc(s *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")
		usr := User{EmailVerificationKey: key}

		err := usr.Read(s.DbWorker)
		if err != nil {
			if err.Error() == "404" {
				service.NotFoundResponse(w, "")
				log.Error(service.DBNotFound)
				return
			}
			log.Error(service.UserReadErr, err)
			service.InternalServerErrorResponse(w, service.UserReadErr, err)
			return
		}

		usr.EmailVerified = true
		err = usr.Update(s.DbWorker)
		if err != nil {
			service.InternalServerErrorResponse(w, service.EmailVerificationErr, err)
			log.Error(service.EmailVerificationErr, err)
			return
		}

		if time.Since(usr.EmailKeyCreatedAt) >= 24*time.Hour {
			service.OkResponse(w, service.DefaultResponse{
				ResultCode: 1,
				HttpCode:   http.StatusGone,
				Messages:   service.VerificationExpired,
				Data:       nil,
			})
			return
		}

		service.OkResponse(w, service.DefaultResponse{
			ResultCode: 0,
			HttpCode:   http.StatusNoContent,
			Messages:   service.VerificationSuccess,
			Data:       nil,
		})
		log.Info(service.VerificationSuccess)
	}
}

// emailResendFunc     godoc
//
//	@Summary		Re-verify email
//	@Description	Re-verify email
//	@Tags			Email
//	@Produce		json
//	@Param			email	path		string					true	"key"
//	@Success		204		{object}	service.DefaultResponse	"No Content"
//	@Failure		400		{object}	service.errorResponse	"Bad request"
//	@Failure		404		{object}	service.errorResponse	"Not Found"
//	@Failure		500		{object}	service.errorResponse	"Internal Server Error"
//	@Router			/reVerifyEmail/{email} [post]
func emailResendFunc(s *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		email := r.PathValue("email")

		usr := User{Email: email}
		err = usr.Read(s.DbWorker)
		if err != nil {
			if err.Error() == "404" {
				service.NotFoundResponse(w, "")
				log.Error(service.DBNotFound)
				return
			}
			log.Error(service.UserReadErr, err)
			service.InternalServerErrorResponse(w, service.UserReadErr, err)
			return
		}

		if usr.EmailVerified {
			service.OkResponse(w, service.DefaultResponse{
				ResultCode: 0,
				HttpCode:   http.StatusNoContent,
				Messages:   service.EmailAlreadyVerified,
				Data:       nil,
			})
		}

		newKey, err := generateEmailVerificationKey()
		if err != nil {
			service.InternalServerErrorResponse(w, service.VerificationKeyErr, err)
			log.Error(service.VerificationKeyErr, err)
			return
		}

		usr.EmailVerificationKey = newKey
		usr.EmailKeyCreatedAt = time.Now()

		err = usr.Update(s.DbWorker)
		if err != nil {
			service.InternalServerErrorResponse(w, service.EmailVerificationErr, err)
			log.Error(service.EmailVerificationErr, err)
			return
		}

		err = sendVerificationEmail(usr.Email, newKey, s)
		if err != nil {
			log.Error(service.EmailSendErr, err)
			return
		}
		log.Debug("Verification link sent on email resend")

		service.OkResponse(w, service.DefaultResponse{
			ResultCode: 0,
			HttpCode:   http.StatusNoContent,
			Messages:   service.VerificationKeySent,
			Data:       nil,
		})
		log.Info(service.VerificationKeySent, " on email resend")
	}
}

func sendVerificationEmail(email string, verificationKey string, s *Service) error {
	m := gomail.NewMessage()

	m.SetHeader("From", s.Config.Config.EmailReply)
	m.SetHeader("To", email)
	m.SetHeader("Subject", service.EmailSubject)
	body := fmt.Sprintf("<html>Please verify your e-mail by following this link <a href='http://localhost:9000/api/v1/verifyEmail/%[1]s'>Verify</a> Key to copypaste: %[1]s</html>", verificationKey)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		s.Config.Config.EmailService,
		25,
		s.Config.Config.EmailLogin,
		s.Config.Config.EmailPass)

	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
