package user

import "errors"

func (u *User) CheckRequiredFields() error {
	if u.Email == "" {
		return errors.New("email is required")
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
