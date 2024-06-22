package user

type LoginForm struct {
	Username string
	Password string
}

func Login(lf *LoginForm) (User, error) {
	user, err := getUser(lf.Username)

	if err != nil {
		return user, nil
	}

	err = comparePasswords(user.Password, lf.Password)
	if err != nil {
		return user, nil
	}
	// TODO swith to cookie
	return user, nil
}

func getUser(u string) (User, error) {
	var user User
	err := DB.Where("username = ?", u).First(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}
