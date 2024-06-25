package todoList

import (
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func Init(d *gorm.DB) {
	DB = d

	err := DB.AutoMigrate(&TodoList{})
	if err != nil {
		log.Fatal(err)
	}

	err = DB.AutoMigrate(&Task{})
	if err != nil {
		log.Fatal(err)
	}
}
