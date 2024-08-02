package user

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model  `json:"-"`
	ID          int       `json:"-"`
	Email       string    `json:"email" binding:"required" example:"example@email.box"  extensions:"x-order=1"`
	Password    string    `json:"password" binding:"required" example:"Very!Strong1Pa$$word"  extensions:"x-order=2"`
	Username    string    `json:"login" extensions:"x-order=3"`
	Name        string    `json:"name" extensions:"x-order=4"`
	Surname     string    `json:"surname" extensions:"x-order=5"`
	UserUUID    uuid.UUID `json:"-"`
	IsSuperuser bool      `json:"-"`
}

type updateUser struct {
	Username string `json:"login" extensions:"x-order=1"`
	Email    string `json:"email" binding:"required" example:"example@email.box" extensions:"x-order=2"`
	Name     string `json:"name" extensions:"x-order=3"`
	Surname  string `json:"surname" extensions:"x-order=4"`
}

type meModel struct {
	UserUUID uuid.UUID `json:"id" extensions:"x-order=1"`
	Email    string    `json:"email" extensions:"x-order=2"`
	Username string    `json:"login" extensions:"x-order=3"`
}

type meResponse struct {
	ResultCode int      `json:"resultCode" extensions:"x-order=1"`
	HttpCode   int      `json:"httpCode" extensions:"x-order=2"`
	Messages   []string `json:"messages" extensions:"x-order=3"`
	Data       meModel  `json:"data" extensions:"x-order=4"`
}

type loginUserModel struct {
	Email    string `json:"email" extensions:"x-order=1"`
	Password string `json:"password" extensions:"x-order=2"`
}

func (l *User) Create() error {
	var err error
	l.UserUUID = uuid.New()
	l.IsSuperuser = false

	l.Password, err = hashPassword(l.Password)
	if err != nil {
		return err
	}

	err = DB.Create(l).Error
	if err != nil {
		return err
	}
	return nil
}

func (l *User) Read() error {
	err := DB.
		Where("user_uuid = ?", l.UserUUID).
		First(l).Error
	if err != nil {
		return err
	}
	return nil
}

func (l *User) Update() error {
	var err error

	err = DB.
		Where("user_uuid = ?", l.UserUUID).
		Updates(l).
		Error
	if err != nil {
		return err
	}

	return nil
}

func (l *User) Delete() error {
	result := DB.
		Where("user_uuid = ?", l.UserUUID).
		Delete(l)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
