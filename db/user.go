package db

import (
	"gorm.io/gorm"
	"log"
)

type User struct {
	gorm.Model
	ID       int
	Nickname string
	Name     string
	Surname  string
	Email    string
	Password string
}

func initUsersTable() {
	err := DB.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}
}

func CreateUser(data []byte) error {
	var err error
	user := new(User)
	err = deserializeJSON(data, user)
	if err != nil {
		return err
	}
	err = DB.Create(&user).Error
	return err
}

func ReadUserByID(id int) ([]byte, error) {
	var user User
	err := DB.First(&user, id).Error

	bytes, err := serializeJSON(user)
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
	err = deserializeJSON(data, &user)
	if err != nil {
		return err
	}
	err = DB.Save(&user).Error
	return err
}

func DeleteUserByID(id int) error {
	err := DB.Delete(&User{}, id).Error
	return err
}
