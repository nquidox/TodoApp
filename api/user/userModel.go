package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"todoApp/types"
)

type User struct {
	gorm.Model  `json:"-"`
	ID          int       `json:"-"`
	Email       string    `json:"email" binding:"required" example:"example@email.box" extensions:"x-order=1"`
	Password    string    `json:"password" binding:"required" example:"Very!Strong1Pa$$word" extensions:"x-order=2"`
	Username    string    `json:"login" extensions:"x-order=3"`
	Name        string    `json:"name" extensions:"x-order=4"`
	Surname     string    `json:"surname" extensions:"x-order=5"`
	UserUUID    uuid.UUID `json:"-"`
	IsSuperuser bool      `json:"-"`
}

type updateUser struct {
	UserUUID uuid.UUID `json:"-"`
	Username string    `json:"login" extensions:"x-order=1"`
	Email    string    `json:"email" binding:"required" example:"example@email.box" extensions:"x-order=2"`
	Name     string    `json:"name" extensions:"x-order=3"`
	Surname  string    `json:"surname" extensions:"x-order=4"`
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

func (u *User) Create(wrk types.DatabaseWorker) error {
	var err error
	u.UserUUID = uuid.New()
	u.IsSuperuser = false

	u.Password, err = hashPassword(u.Password)
	if err != nil {
		return err
	}

	err = wrk.CreateRecord(u)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Read(wrk types.DatabaseWorker) error {
	params := map[string]any{"user_uuid": u.UserUUID}
	err := wrk.ReadOneRecord(u, params)
	if err != nil {
		return err
	}
	return nil
}

func (u *updateUser) Update(wrk types.DatabaseWorker) error {
	params := map[string]any{"user_uuid": u.UserUUID}
	err := wrk.UpdateRecord(u, params)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Delete(wrk types.DatabaseWorker) error {
	params := map[string]any{"user_uuid": u.UserUUID}
	err := wrk.DeleteRecord(u, params)
	if err != nil {
		return err
	}
	return nil
}

func (m *meModel) Read(wrk types.DatabaseWorker) error {
	var usr User
	params := map[string]any{"user_uuid": m.UserUUID}
	err := wrk.ReadOneRecord(&usr, params)

	m.Email = usr.Email
	m.Username = usr.Username

	if err != nil {
		return err
	}
	return nil
}
