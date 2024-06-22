package user

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitUsers() {
	err := DB.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}

	err = DB.AutoMigrate(&Session{})
	if err != nil {
		log.Fatal(err)
	}
}

func serializeJSON(v interface{}) ([]byte, error) {
	marshal, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}

func deserializeJSON(data []byte, s interface{}) error {
	err := json.Unmarshal(data, &s)
	return err
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func comparePasswords(hashedPassword string, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
