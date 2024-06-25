package user

import (
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func Init(d *gorm.DB) {
	DB = d
	err := DB.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}

	err = DB.AutoMigrate(&Session{})
	if err != nil {
		log.Fatal(err)
	}
}
