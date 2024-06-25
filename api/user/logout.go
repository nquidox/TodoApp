package user

func Logout(cookieToken string) error {
	err := DropSession(cookieToken)
	if err != nil {
		return err
	}

	return nil
}
