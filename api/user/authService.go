package user

import "todoApp/types"

type AuthService struct {
	types.AuthUser
}

func (a *AuthService) IsUserLoggedIn(dbw types.DatabaseWorker, tokenValue string) (types.AuthUser, error) {
	s := Session{Token: tokenValue}
	err := s.Read(dbw)
	if err != nil {
		return a.AuthUser, err
	}

	params := map[string]any{"user_uuid": s.UserUuid}

	err = dbw.ReadRecordSubmodel(User{}, a, params)
	if err != nil {
		return a.AuthUser, err
	}
	return a.AuthUser, nil
}
