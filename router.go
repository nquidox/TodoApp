package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"todoApp/user"
)

type ApiServer struct {
	Addr string
}

func NewApiServer(host, port string) *ApiServer {
	return &ApiServer{Addr: host + ":" + port}
}

func (s *ApiServer) Run() error {
	router := http.NewServeMux()
	server := &http.Server{Addr: s.Addr, Handler: router}

	router.HandleFunc("GET /api/v1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "API v1 is ready.")
	})

	router.HandleFunc("POST /api/v1/create", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		err = user.CreateUser(body)

		if err != nil {
			fmt.Fprint(w, "Error creating user", err)
		} else {
			fmt.Fprint(w, "User created")
		}

	})

	router.HandleFunc(fmt.Sprintf("GET /api/v1/read/{id}"), func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		usr, err := user.ReadUserByID(userID)
		if err != nil {
			fmt.Fprint(w, "Error reading user", http.StatusBadRequest)
		} else {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, string(usr))
		}
	})

	router.HandleFunc("UPDATE /api/v1/update/{id}", func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.Atoi(r.PathValue("id"))
		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
		}
		err = user.UpdateUserByID(userID, body)
		if err != nil {
			fmt.Fprint(w, "Error updating user", http.StatusBadRequest)
		}
	})

	router.HandleFunc("DELETE /api/v1/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
		}

		err = user.DeleteUserByID(userID)
		if err != nil {
			fmt.Fprint(w, "Error deleting user", http.StatusBadRequest)
		}
	})

	router.HandleFunc("POST /api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
		}

		cookie, err := user.Login(&user.LoginForm{
			Username: r.Form.Get("username"),
			Password: r.Form.Get("password"),
		})

		if err != nil {
			fmt.Fprint(w, "Error logging in", err, http.StatusBadRequest)
		} else {

			w.Header().Set("Content-Type", "application/json")
			http.SetCookie(w, &cookie)
		}
	})

	log.Printf("Starting server on %s", s.Addr)
	return server.ListenAndServe()
}
