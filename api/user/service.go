package user

import (
	"crypto/rand"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func comparePasswords(hashedPassword string, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

func generateEmailVerificationKey() (string, error) {
	bytes := make([]byte, 16)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	value := hex.EncodeToString(bytes)

	log.Debug("Generated email verification key")

	return value, nil
}
