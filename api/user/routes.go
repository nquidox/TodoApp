package user

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func AddRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /api/v1/create", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		err = CreateUser(body)
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

		usr, err := ReadUserByID(userID)
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
		err = UpdateUserByID(userID, body)
		if err != nil {
			fmt.Fprint(w, "Error updating user", http.StatusBadRequest)
		}
	})

	router.HandleFunc("DELETE /api/v1/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
		}

		err = DeleteUserByID(userID)
		if err != nil {
			fmt.Fprint(w, "Error deleting user", http.StatusBadRequest)
		}
	})

	router.HandleFunc("POST /api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
		}

		cookie, err := Login(&LoginForm{
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

	router.HandleFunc("GET /api/v1/logout", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			fmt.Fprint(w, "Error getting cookie", http.StatusBadRequest)
		}

		err = Logout(cookie.Value)
		if err != nil {
			fmt.Fprint(w, "Error in attempt to logout", http.StatusBadRequest)
		}
	})
}
