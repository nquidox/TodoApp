package user

import (
	"errors"
	"net/mail"
)

func (u *User) CheckRequiredFields() error {
	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		return err
	}

	if u.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func (l *loginUserModel) CheckRequiredFields() error {
	if l.Email == "" {
		return errors.New("email is required")
	}

	if l.Password == "" {
		return errors.New("password is required")
	}
	return nil
}
