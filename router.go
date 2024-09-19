package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"todoApp/config"
	_ "todoApp/docs"
)

type ApiServer struct {
	Addr string
}

func NewApiServer(host, port string) *ApiServer {
	return &ApiServer{Addr: host + ":" + port}
}

func (s *ApiServer) Run(router *http.ServeMux, cors *config.CORSConfig) error {
	server := &http.Server{Addr: s.Addr, Handler: corsMiddleware(router, cors)}

	router.HandleFunc("GET /api/v1", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprint(w, "API v1 is ready.")
		if err != nil {
			log.Fatal(err)
			return
		}
	})

	router.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	log.Info("Starting server on ", s.Addr)
	return server.ListenAndServe()
}

func corsMiddleware(next http.Handler, cors *config.CORSConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := cors.AllowedOrigins
		origin := r.Header.Get("Origin")

		for _, o := range allowedOrigins {
			if o == origin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
