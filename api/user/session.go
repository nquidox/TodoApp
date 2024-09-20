package user

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"time"
	"todoApp/api/service"
)

type Session struct {
	gorm.Model
	Id         int
	UserUuid   uuid.UUID
	Token      string
	ClientInfo string
	Expires    time.Time
}

func (s *Session) Create(wrk dbWorker, userUuid uuid.UUID, salt []byte, userAgent string) (http.Cookie, error) {
	token, err := generateToken(32, salt)
	if err != nil {
		return http.Cookie{}, err
	}

	expires := time.Now().Add(3 * time.Hour * 24 * 365)

	s.UserUuid = userUuid
	s.Token = token
	s.ClientInfo = userAgent
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

func (s *Session) Read(wrk dbWorker) error {
	params := map[string]any{service.SessionTokenName: s.Token}
	err := wrk.ReadOneRecord(&s, params)
	if err != nil {
		return err
	}
	return nil
}

func (s *Session) Delete(wrk dbWorker) error {
	params := map[string]any{service.SessionTokenName: s.Token}
	err := wrk.ReadOneRecord(s, params)
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
		Name:        service.SessionTokenName,
		Value:       token,
		Path:        "/",
		Expires:     expires,
		Secure:      true,
		HttpOnly:    false,
		SameSite:    http.SameSiteNoneMode,
		Partitioned: true,
	}
	log.WithFields(log.Fields{
		"name":  service.SessionTokenName,
		"value": token,
	}).Debug()
	return cookie, nil
}

func generateToken(length int, salt []byte) (string, error) {
	bytes := make([]byte, length)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	bytes = append(bytes, salt...)
	value := hex.EncodeToString(bytes)

	log.WithFields(log.Fields{
		"value": value,
	}).Debug("Generated token")

	return value, nil
}
