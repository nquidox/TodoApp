package db

import (
	"gorm.io/gorm"
	"log"
)

type User struct {
	gorm.Model
	ID      int
	Name    string
	Surname string
	Email   string
}

func initUsersTable() {
	err := DB.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}
}

func CreateUser(user User) error {
	err := DB.Create(&user).Error
	return err
}
