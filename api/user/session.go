package user

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

var SALT []byte

type Session struct {
	gorm.Model
	Id         int
	UserUuid   uuid.UUID
	Token      string
	ClientInfo string
	Expires    time.Time
}

func getTokenValue(r *http.Request) (string, error) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Println(err)
		return "", err
	}
	return token.Value, nil
}

func tokenIsValid(token string) (bool, uuid.UUID) {
	var s Session
	result := DB.Where("token = ?", token).First(&s)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, uuid.Nil
	}

	if result.Error != nil {
		log.Println(result.Error)
		return false, uuid.Nil
	}

	return true, s.UserUuid
}

func createSession(u uuid.UUID) (http.Cookie, error) {
	token, err := generateToken(32)
	if err != nil {
		return http.Cookie{}, err
	}

	expires := time.Now().Add(3 * time.Hour * 24 * 365)

	session := Session{
		UserUuid:   u,
		Token:      token,
		ClientInfo: "",
		Expires:    expires,
	}

	err = DB.Create(&session).Error
	if err != nil {
		return http.Cookie{}, err
	}

	cookie, err := createSessionCookie(token, expires)
	if err != nil {
		return http.Cookie{}, err
	}

	return cookie, nil
}

func dropSession(cookieToken string) error {
	var session Session
	err := DB.Where("token = ?", cookieToken).First(&session).Error
	if err != nil {
		return err
	}
	err = DB.Delete(&session).Error
	if err != nil {
		return err
	}
	return nil
}

func createSessionCookie(token string, expires time.Time) (http.Cookie, error) {
	cookie := http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  expires,
		Secure:   false,
		HttpOnly: true,
	}
	return cookie, nil
}

func generateToken(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	bytes = append(bytes, SALT...)
	return hex.EncodeToString(bytes), nil
}
