package infoPages

func addRoutes(s *Service) {
	s.Router.HandleFunc("GET /info", infoPageHandler)
}
