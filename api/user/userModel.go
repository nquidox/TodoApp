package user

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       int
	Username string `json:"username"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Uuid     uuid.UUID
}

func (u *User) Create() error {
	var err error
	u.Uuid = uuid.New()

	u.Password, err = hashPassword(u.Password)
	if err != nil {
		return err
	}

	err = DB.Create(u).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Read() error {
	err := DB.
		Where("uuid = ?", u.Uuid).
		First(u).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Update() error {
	var err error

	err = DB.
		Where("uuid = ?", u.Uuid).
		Updates(u).
		Error
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Delete() error {
	result := DB.
		Where("uuid = ?", u.Uuid).
		Delete(u)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
