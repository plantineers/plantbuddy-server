// Author: Maximilian Floto
package auth

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/utils"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HttpMethodNotAllowedResponse(w, "Method not allowed. Login only supports GET requests.")
		return
	}
	handleLoginGet(w, r)
}

func handleLoginGet(w http.ResponseWriter, r *http.Request) {
	safeUser, err := authUser(r)
	switch err {
	case ErrWrongCredentials:
		utils.HttpForbiddenResponse(w, "Wrong credentials")
	case ErrNoCredentials:
		utils.HttpBadRequestResponse(w, "No credentials supplied")
	case nil:
		b, err := json.Marshal(safeUser)
		if err != nil {
			msg := fmt.Sprintf("Error converting user %s to JSON: %s", safeUser.Name, err.Error())
			utils.HttpInternalServerErrorResponse(w, msg)
			return
		}

		log.Printf("User %s logged in", safeUser.Name)
		utils.HttpOkResponse(w, b)
	default:
		msg := fmt.Sprintf("Error while authenticating user: %s", err.Error())
		utils.HttpInternalServerErrorResponse(w, msg)
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
