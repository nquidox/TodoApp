package todoList

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"todoApp/api/service"
)

func (a *authUser) isAuth(w http.ResponseWriter, r *http.Request, s *Service) error {
	token, err := r.Cookie(service.SessionTokenName)
	if err != nil {
		service.UnauthorizedResponse(w, "")
		log.Error(service.TokenReadErr, err.Error())
		return err
	}

	authUsr, err := s.AuthWorker.IsUserLoggedIn(s.DbWorker, token.Value)
	if err != nil {
		service.UnauthorizedResponse(w, "")
		log.Error(service.AuthErr, err)
		return err
	}

	a.AuthUser = authUsr
	return nil
}
