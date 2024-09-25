package infoPages

import "net/http"

type Service struct {
	Router *http.ServeMux
}

func Init(s *Service) {
	addRoutes(s)
}
