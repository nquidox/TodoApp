package user

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model  `json:"-"`
	ID          int       `json:"-"`
	Username    string    `json:"login"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	Email       string    `json:"email" binding:"required" example:"example@email.box"`
	Password    string    `json:"password" binding:"required" example:"Very!Strong1Pa$$word"`
	Uuid        uuid.UUID `json:"-"`
	IsSuperuser bool      `json:"-"`
}

func (u *User) Create() error {
	var err error
	u.Uuid = uuid.New()
	u.IsSuperuser = false

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
