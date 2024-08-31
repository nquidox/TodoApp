package user

import (
	"errors"
	"net/mail"
)

func (u *User) CheckRequiredFields() error {
	err := validateRequiredFields(u.Email, u.Password)
	if err != nil {
		return err
	}
	return nil
}

func (l *loginUserModel) CheckRequiredFields() error {
	err := validateRequiredFields(l.Email, l.Password)
	if err != nil {
		return err
	}
	return nil
}

func validateRequiredFields(email, password string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return err
	}

	if password == "" {
		return errors.New("password is required")
	}
	return nil
}
