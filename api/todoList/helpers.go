package todoList

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"todoApp/api/service"
)

func (a *authUser) isAuth(w http.ResponseWriter, r *http.Request) error {
	token, err := r.Cookie("token")
	if err != nil {
		service.UnauthorizedResponse(w, "")
		log.Error(service.TokenReadErr, err.Error())
		return err
	}

	authUsr, err := aw.IsUserLoggedIn(dbw, token.Value)
	if err != nil {
		service.UnauthorizedResponse(w, "")
		log.Error(service.AuthErr, err)
		return err
	}

	a.UserUUID = authUsr.UserUUID
	a.IsSuperuser = authUsr.IsSuperuser

	return nil
}
