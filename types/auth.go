package types

import (
	"github.com/google/uuid"
)

type AuthWorker interface {
	IsUserLoggedIn(wrk DatabaseWorker, tokenValue string) (AuthUser, error)
}

type AuthUser struct {
	UserUUID    uuid.UUID
	IsSuperuser bool
}
