// Author: Maximilian Floto
package auth

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/plantineers/plantbuddy-server/db"
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
	safeUser, err := authUser(r)
	switch err {
	case ErrWrongCredentials:
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Wrong credentials"))
	case ErrNoCredentials:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No credentials supplied"))
	case nil:
		b, err := json.Marshal(safeUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func getUserByName(name string) (*User, error) {
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

func authUser(r *http.Request) (*SafeUser, error) {
	authHeader := r.Header.Get("Authorization")
	disasembledAuthHeader := strings.Split(authHeader, " ")
	if len(disasembledAuthHeader) != 2 || disasembledAuthHeader[0] != "Basic" {
		return nil, ErrInvalidAuthHeader
	}

	decodedAuthHeaderDigest := make([]byte, base64.StdEncoding.DecodedLen(len(disasembledAuthHeader[1])))
	n, err := base64.StdEncoding.Decode(decodedAuthHeaderDigest, []byte(disasembledAuthHeader[1]))
	if err != nil {
		return nil, ErrInvalidAuthHeader
	}

	decodedAuthHeader := strings.Split(string(decodedAuthHeaderDigest[:n]), ":")
	userName := decodedAuthHeader[0]
	password := decodedAuthHeader[1]

	// Get user from db
	user, err := getUserByName(userName)
	if err == sql.ErrNoRows { // User not found
		return nil, ErrWrongCredentials
	} else if err != nil { // Database error
		return nil, err
	}

	// Check password
	password = utils.HashPassword(password)
	if password != user.Password {
		return nil, ErrWrongCredentials
	}

	return &SafeUser{
		Id:   user.Id,
		Name: user.Name,
		Role: user.Role,
	}, nil
}
