package user

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

func Authorized(r *http.Request, userId uuid.UUID) error {
	var userSession Session
	token, err := r.Cookie("token")
	if err != nil {
		return err
	}

	err = DB.Where("token = ?", token.Value).First(&userSession).Error
	if err != nil {
		return err
	}

	if userSession.UserUuid != userId {
		return errors.New("invalid token")
	}

	return nil
}

// This function gets user's UUID by session token. Also ensures that user is authorized.
func getUserUUIDFromToken(token string) (uuid.UUID, error) {
	var userSession Session
	result := DB.Where("token = ?", token).First(&userSession)

	if result.Error != nil {
		return uuid.Nil, result.Error
	}

	if result.RowsAffected == 0 {
		return uuid.Nil, errors.New("no such user")
	}

	return userSession.UserUuid, nil
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
