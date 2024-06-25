package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"todoApp/service"
)

type User struct {
	gorm.Model
	ID       int
	Username string
	Name     string
	Surname  string
	Email    string
	Password string
	Uuid     uuid.UUID
}

func CreateUser(data []byte) error {
	var err error
	user := new(User)

	err = service.DeserializeJSON(data, user)
	if err != nil {
		return err
	}

	user.Password, err = hashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Uuid = uuid.New()

	err = DB.Create(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func ReadUserByID(id int) ([]byte, error) {
	var user User
	err := DB.First(&user, id).Error

	bytes, err := service.SerializeJSON(user)
	if err != nil {
		return nil, err
	}

	return bytes, err
}

func UpdateUserByID(id int, data []byte) error {
	var user User
	err := DB.First(&user, id).Error
	if err != nil {
		return err
	}
	err = service.DeserializeJSON(data, &user)
	if err != nil {
		return err
	}
	err = DB.Save(&user).Error
	return err
}

func DeleteUserByID(id int) error {
	err := DB.Delete(&User{}, id).Error
	if err != nil {
		return err
	}
	return nil
}
