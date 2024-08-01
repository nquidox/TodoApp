package user

import "errors"

func (l *User) CheckRequiredFields() error {
	if l.Email == "" {
		return errors.New("email is required")
	}

	if l.Password == "" {
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
