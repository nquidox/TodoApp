package user

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type Session struct {
	Id         int
	Uuid       uuid.UUID
	Token      string
	ClientInfo string
	Expires    time.Time
}

func CreateSession(u uuid.UUID) (http.Cookie, error) {
	token, err := generateToken(32)
	if err != nil {
		return http.Cookie{}, err
	}

	expires := time.Now().Add(3 * time.Hour * 24 * 365)

	session := Session{
		Uuid:       u,
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
	return hex.EncodeToString(bytes), nil
}
