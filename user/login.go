package user

import "net/http"

type LoginForm struct {
	Username string
	Password string
}

func Login(lf *LoginForm) (http.Cookie, error) {
	usr, err := getUser(lf.Username)
	if err != nil {
		return http.Cookie{}, err
	}

	err = comparePasswords(usr.Password, lf.Password)
	if err != nil {
		return http.Cookie{}, err
	}

	session, err := CreateSession(usr.Uuid)
	if err != nil {
		return http.Cookie{}, err
	}

	return session, nil
}

func getUser(u string) (User, error) {
	var user User
	err := DB.Where("username = ?", u).First(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}
