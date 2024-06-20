package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"todoApp/config"
)

type User struct {
	gorm.Model
	ID      int
	Name    string
	Surname string
	Email   string
}

var DB *gorm.DB

func ConnectToDB(c *config.Config) {
	var err error
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Config.Host, c.Config.Port, c.Config.User, c.Config.Password, c.Config.Dbname, c.Config.Sslmode)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal(err)
	}

	err = DB.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}
}

func CreateUser(user User) error {
	err := DB.Create(&user).Error
	return err
}
