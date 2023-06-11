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

// LoginHandler handles login requests.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// LoginHandler only accepts GET requests.
	if r.Method != http.MethodGet {
		utils.HttpMethodNotAllowedResponse(w, "Allowed methods: GET")
		return
	}
	handleLoginGet(w, r)
}

// handleLoginGet handles GET requests to the login endpoint.
func handleLoginGet(w http.ResponseWriter, r *http.Request) {
	safeUser, err := authUser(r)
	switch err {
	case ErrWrongCredentials:
		utils.HttpForbiddenResponse(w, "Wrong credentials")
	case ErrNoCredentials:
		utils.HttpBadRequestResponse(w, "No credentials supplied")
	case nil:
		// User is authenticated, return SafeUser object
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

// getUserByName returns a user from the database by name.
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

// authUser authorizes a user by checking the Authorization header.
// It follows the HTTP Basic Auth scheme.
func authUser(r *http.Request) (*SafeUser, error) {
	// Get Authorization header value
	authHeader := r.Header.Get("Authorization")
	disassembledAuthHeader := strings.Split(authHeader, " ")
	if len(disassembledAuthHeader) != 2 || disassembledAuthHeader[0] != "Basic" {
		return nil, ErrInvalidAuthHeader
	}

	// Decode header value
	decodedAuthHeaderDigest := make([]byte, base64.StdEncoding.DecodedLen(len(disassembledAuthHeader[1])))
	n, err := base64.StdEncoding.Decode(decodedAuthHeaderDigest, []byte(disassembledAuthHeader[1]))
	if err != nil {
		return nil, ErrInvalidAuthHeader
	}

	// Extract username and password from decoded header value
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
