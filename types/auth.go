package types

import (
	"github.com/google/uuid"
)

type AuthWorker interface {
	IsUserLoggedIn(wrk DatabaseWorker, id uuid.UUID) (AuthUser, error)
}

type AuthUser struct {
	UserUUID    uuid.UUID
	IsSuperuser bool
}
