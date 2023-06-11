// Author: Maximilian Floto
package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/plantineers/plantbuddy-server/utils"

	"github.com/plantineers/plantbuddy-server/db"
)

// UsersHandler handles all requests to the users endpoint.
func UsersHandler(w http.ResponseWriter, r *http.Request) {
	// UsersHandler only accepts GET requests.
	if r.Method != http.MethodGet {
		utils.HttpMethodNotAllowedResponse(w, "Allowed methods: GET")
		return
	}
	handleUsersGet(w, r)
}

// handleUsersGet handles GET requests to the users endpoint.
func handleUsersGet(w http.ResponseWriter, r *http.Request) {
	users, err := getAllUsers()
	switch err {
	case nil:
		b, err := json.Marshal(users)
		if err != nil {
			msg := fmt.Sprintf("Error converting users to JSON: %s", err.Error())
			utils.HttpInternalServerErrorResponse(w, msg)
			return
		}
		log.Printf("Loaded %d users", len(users))
		utils.HttpOkResponse(w, b)
	default:
		msg := fmt.Sprintf("Error while loading users: %s", err.Error())
		utils.HttpInternalServerErrorResponse(w, msg)
	}
}

// getAllUsers returns all user ids from the database.
func getAllUsers() ([]int64, error) {
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

	return repo.GetAll()
}
