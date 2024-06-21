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
	user := new(User)
	deserializeJSON(data, user)
	err := DB.Create(&user).Error
	return err
}
