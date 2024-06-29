package user

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	err = CreateUser(body)
	if err != nil {
		fmt.Fprint(w, "Error creating user", err)
	} else {
		fmt.Fprint(w, "User created")
	}
}

func ReadUserHandler(w http.ResponseWriter, r *http.Request) {
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
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("id"))
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
	}
	err = UpdateUserByID(userID, body)
	if err != nil {
		fmt.Fprint(w, "Error updating user", http.StatusBadRequest)
	}
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
	}

	err = DeleteUserByID(userID)
	if err != nil {
		fmt.Fprint(w, "Error deleting user", http.StatusBadRequest)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
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
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		fmt.Fprint(w, "Error getting cookie", http.StatusBadRequest)
	}

	err = Logout(cookie.Value)
	if err != nil {
		fmt.Fprint(w, "Error in attempt to logout", http.StatusBadRequest)
	}
}
