package user

import (
	"crypto/tls"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"net/http"
	"todoApp/api/service"
)

// emailFunc     godoc
//
//	@Summary		Verify email
//	@Description	Verify email by key
//	@Tags			Auth
//	@Produce		json
//	@Param			key	path		string					true	"key"
//	@Success		204	{object}	service.DefaultResponse	"No Content"
//	@Failure		400	{object}	service.errorResponse	"Bad request"
//	@Failure		404	{object}	service.errorResponse	"Not Found"
//	@Failure		500	{object}	service.errorResponse	"Internal Server Error"
//	@Router			/verifyEmail/{key} [post]
func emailFunc(s *Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		log.Info("User email verified")
	})
}

func sendVerificationEmail(email string, verficationKey string, s *Service) error {
	m := gomail.NewMessage()

	m.SetHeader("From", s.Config.Config.EmailReply)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "")
	body := fmt.Sprintf("<html>Please verify your e-mail by following this link %s</html>", verficationKey)
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
