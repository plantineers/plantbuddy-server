// Author: Maximilian Floto
package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/utils"
)

// UserCreateHandler handles user creation requests.
func UserCreateHandler(w http.ResponseWriter, r *http.Request) {
	// UserCreateHandler only accepts POST requests.
	if r.Method != http.MethodPost {
		utils.HttpMethodNotAllowedResponse(w, "Allowed methods: POST")
		return
	}
	handleUserPost(w, r)
}

// UserHandler handles all requests to the user endpoint, except POST.
func UserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := utils.PathParameterFilter(r.URL.Path, "/v1/user/")
	if err != nil {
		utils.HttpBadRequestResponse(w, "No user id supplied")
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleUserGet(w, r, id)
	case http.MethodPut:
		handleUserPut(w, r, id)
	case http.MethodDelete:
		handleUserDelete(w, r, id)
	default:
		utils.HttpMethodNotAllowedResponse(w, "Allowed methods: GET, PUT, DELETE")
	}
}

// handleUserPost handles POST requests to the user endpoint.
func handleUserPost(w http.ResponseWriter, r *http.Request) {
	var user User
	// Get user from request body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		msg := fmt.Sprintf("Error decoding new user: %s", err.Error())
		utils.HttpBadRequestResponse(w, msg)
		return
	}

	user.Password = utils.HashPassword(user.Password)

	createdUser, err := createUser(&user)
	switch err {
	case nil:
		safeUser := &SafeUser{
			Id:   createdUser.Id,
			Name: createdUser.Name,
			Role: createdUser.Role,
		}

		b, err := json.Marshal(safeUser)
		if err != nil {
			msg := fmt.Sprintf("Error converting safe user %s to JSON: %s", safeUser.Name, err.Error())
			utils.HttpInternalServerErrorResponse(w, msg)
			return
		}

		msg := fmt.Sprintf("Created user %s", safeUser.Name)
		location := fmt.Sprintf("/v1/user/%d", safeUser.Id)
		utils.HttpCreatedResponse(w, b, location, msg)
	case ErrUserAlreadyExists:
		msg := fmt.Sprintf("User %s already exists", user.Name)
		utils.HttpConflictResponse(w, msg)
	default:
		msg := fmt.Sprintf("Error while creating user %s", user.Name)
		utils.HttpInternalServerErrorResponse(w, msg)
	}
}

// handleUserGet handles GET requests to the user endpoint.
func handleUserGet(w http.ResponseWriter, r *http.Request, id int64) {
	user, err := getUserById(id)
	switch err {
	case nil:
		safeUser := &SafeUser{
			Id:   user.Id,
			Name: user.Name,
			Role: user.Role,
		}
		b, err := json.Marshal(safeUser)
		if err != nil {
			msg := fmt.Sprintf("Error converting safe user %s to JSON: %s", safeUser.Name, err.Error())
			utils.HttpInternalServerErrorResponse(w, msg)
			return
		}

		log.Printf("User %s loaded", safeUser.Name)
		utils.HttpOkResponse(w, b)
	case sql.ErrNoRows:
		msg := fmt.Sprintf("User with id %d not found", id)
		utils.HttpNotFoundResponse(w, msg)
	default:
		msg := fmt.Sprintf("Error while getting user with id %d: %s", id, err.Error())
		utils.HttpBadRequestResponse(w, msg)
	}
}

// handleUserPut handles PUT requests to the user endpoint.
func handleUserPut(w http.ResponseWriter, r *http.Request, id int64) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		msg := fmt.Sprintf("Error decoding new user: %s", err.Error())
		utils.HttpBadRequestResponse(w, msg)
		return
	}

	user.Password = utils.HashPassword(user.Password)
	user.Id = id

	err = updateUser(&user)
	if err != nil {
		msg := fmt.Sprintf("Error converting user %s to JSON: %s", user.Name, err.Error())
		utils.HttpInternalServerErrorResponse(w, msg)
		return
	}

	safeUser := &SafeUser{
		Id:   user.Id,
		Name: user.Name,
		Role: user.Role,
	}

	b, err := json.Marshal(safeUser)
	if err != nil {
		msg := fmt.Sprintf("Error converting user %s to JSON: %s", user.Name, err.Error())
		utils.HttpInternalServerErrorResponse(w, msg)
		return
	}

	log.Printf("Updated user %s with id %d", user.Name, user.Id)
	utils.HttpOkResponse(w, b)
}

// handleUserDelete handles DELETE requests to the user endpoint.
func handleUserDelete(w http.ResponseWriter, r *http.Request, id int64) {
	err := deleteUserById(id)

	switch err {
	case nil:
		log.Printf("Deleted user with id %d", id)
		utils.HttpOkResponse(w, nil)
	case ErrCannotDeleteRoot:
		msg := fmt.Sprintf("Error while deleting user with id %d (user is root): %s", id, err.Error())
		utils.HttpForbiddenResponse(w, msg)
	default:
		msg := fmt.Sprintf("Error while deleting user with id %d: %s", id, err.Error())
		utils.HttpInternalServerErrorResponse(w, msg)
	}
}

// createUser creates a new user in the database.
func createUser(user *User) (*User, error) {
	session := db.NewSession()
	defer session.Close()

	err := session.Open()
	if err != nil {
		return nil, err
	}

	repo, err := NewUserRepository(session)
	if err != nil {
		return nil, err
	}

	_, err = repo.GetByName(user.Name)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}

	err = repo.Create(user)
	if err != nil {
		return nil, err
	}

	createdUser, err := repo.GetByName(user.Name)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

// getUserById gets a user from the database by id.
func getUserById(id int64) (*User, error) {
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

	return repo.GetById(id)
}

// updateUser updates a user in the database.
func updateUser(user *User) error {
	session := db.NewSession()
	defer session.Close()

	err := session.Open()
	if err != nil {
		return err
	}

	repo, err := NewUserRepository(session)
	if err != nil {
		return err
	}

	return repo.Update(user)
}

// deleteUserById deletes a user from the database by id.
func deleteUserById(id int64) error {
	session := db.NewSession()
	defer session.Close()

	err := session.Open()
	if err != nil {
		return err
	}

	repo, err := NewUserRepository(session)
	if err != nil {
		return err
	}

	user, err := getUserById(id)
	if err != nil {
		return err
	}

	// Prevent admins from deleting the root user
	if user.Name == "root" {
		return ErrCannotDeleteRoot
	}

	return repo.DeleteById(id)
}
