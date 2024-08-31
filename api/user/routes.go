package user

func addRoutes(s *Service) {
	createUserHandler := createUserFunc(s)
	s.Router.HandleFunc("POST /api/v1/user", createUserHandler)

	readUserHandler := readUserFunc(s)
	s.Router.HandleFunc("GET /api/v1/user/{id}", readUserHandler)

	updateUserHandler := updateUserFunc(s)
	s.Router.HandleFunc("PUT /api/v1/user/{id}", updateUserHandler)

	deleteUserHandler := deleteUserFunc(s)
	s.Router.HandleFunc("DELETE /api/v1/user/{id}", deleteUserHandler)

	loginHandler := loginFunc(s)
	s.Router.HandleFunc("POST /api/v1/login", loginHandler)

	logoutHandler := logoutFunc(s)
	s.Router.HandleFunc("GET /api/v1/logout", logoutHandler)

	meHandler := meFunc(s)
	s.Router.HandleFunc("GET /api/v1/me", meHandler)

	emailHandler := emailFunc(s)
	s.Router.HandleFunc("GET /api/v1/verifyEmail/{key}", emailHandler)

	emailResendHandler := emailResendFunc(s)
	s.Router.HandleFunc("POST /api/v1/reVerifyEmail/{email}", emailResendHandler)
}
