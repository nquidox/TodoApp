package user

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"todoApp/types"

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

func (s *Session) Create(wrk types.DatabaseWorker, userUuid uuid.UUID) (http.Cookie, error) {
	token, err := generateToken(32)
	if err != nil {
		return http.Cookie{}, err
	}

	expires := time.Now().Add(3 * time.Hour * 24 * 365)

	s.UserUuid = userUuid
	s.Token = token
	s.ClientInfo = ""
	s.Expires = expires

	err = wrk.CreateRecord(s)
	if err != nil {
		return http.Cookie{}, err
	}

	cookie, err := createSessionCookie(token, expires)
	if err != nil {
		return http.Cookie{}, err
	}

	return cookie, nil
}

func (s *Session) Read(wrk types.DatabaseWorker) error {
	params := map[string]any{"token": s.Token}
	err := wrk.ReadOneRecord(&s, params)
	if err != nil {
		return err
	}
	return nil
}

func (s *Session) Delete(wrk types.DatabaseWorker) error {
	params := map[string]any{"token": s.Token}
	err := Worker.ReadOneRecord(s, params)
	if err != nil {
		return err
	}

	err = wrk.DeleteRecord(s, params)
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
