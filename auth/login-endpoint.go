// Author: Maximilian Floto
package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/model"
	"github.com/plantineers/plantbuddy-server/utils"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed. Login only supports GET requests."))
		return
	}
	handleLoginGet(w, r)
}

func handleLoginGet(w http.ResponseWriter, r *http.Request) {
	userName := r.Header.Get("X-User-Name")
	password := r.Header.Get("X-User-Password")

	// Check if credentials were supplied
	if userName == "" || password == "" {
		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte("No credentials supplied!"))
		if err != nil {
			return
		}
		return
	}

	// Get user from db
	user, err := getUserByName(userName)
	if err == sql.ErrNoRows { // User not found
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("User not found"))
		return
	} else if err != nil { // Database error
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Check password
	password = utils.HashPassword(password)
	if password != user.Password {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Invalid password"))
		return
	}

	safeUser := &model.SafeUser{
		Id:   user.Id,
		Name: user.Name,
		Role: user.Role,
	}

	b, err := json.Marshal(safeUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func getUserByName(name string) (*model.User, error) {
	var session = db.NewSession()
	defer session.Close()

	err := session.Open()
	if err != nil {
		return nil, err
	}

	repo, err := NewUserRepository(session)
	if err != nil {
		return nil, err
	}

	return repo.GetByName(name)
}
