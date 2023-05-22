package user_management

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/plantineers/plantbuddy-server/model"
	"net/http"
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
	password = fmt.Sprintf("%x", sha256.Sum256([]byte(password+"plantbuddy_salt")))
	if password != user.Password {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Invalid password"))
		return
	}

	safeUser := &model.SafeUser{
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
