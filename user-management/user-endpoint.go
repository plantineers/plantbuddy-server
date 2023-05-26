// Author: Maximilian Floto
package user_management

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/model"
	"github.com/plantineers/plantbuddy-server/utils"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleUserGet(w, r)
	case http.MethodPost:
		handleUserPost(w, r)
	case http.MethodDelete:
		handleUserDelete(w, r)
	case http.MethodPut:
		handleUserPut(w, r)
	}
}

func handleUserGet(w http.ResponseWriter, r *http.Request) {
	id, err := utils.PathParameterFilter(r.URL.Path, "/v1/user/")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	user, err := getUserById(id)
	switch err {
	case nil:
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
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	case sql.ErrNoRows:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
}

func getUserById(id int64) (*model.User, error) {
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

func handleUserPost(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.Password = utils.HashPassword(user.Password)

	createdUser, err := createUser(&user)
	switch err {
	case nil:
		safeUser := &model.SafeUser{
			Id:   createdUser.Id,
			Name: createdUser.Name,
			Role: createdUser.Role,
		}

		b, err := json.Marshal(safeUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(b)
	case ErrUserAlreadyExists:
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error()))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func createUser(user *model.User) (*model.User, error) {
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

// TODO: Move to errors.go
var ErrUserAlreadyExists = errors.New("User already exists")

func handleUserDelete(w http.ResponseWriter, r *http.Request) {
	id, err := utils.PathParameterFilter(r.URL.Path, "/v1/user/")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = deleteUserById(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
	case ErrCannotDeleteAdmin:
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
}

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

	// Prevent admins from deleting other admins
	if user.Role == model.Admin {
		return ErrCannotDeleteAdmin
	}

	return repo.DeleteById(id)
}

// TODO: Move to errors.go
var ErrCannotDeleteAdmin = errors.New("Cannot delete admin user")

func handleUserPut(w http.ResponseWriter, r *http.Request) {
	id, err := utils.PathParameterFilter(r.URL.Path, "/v1/user/")
	var user model.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.Password = utils.HashPassword(user.Password)
	user.Id = id

	err = updateUser(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
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

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func updateUser(user *model.User) error {
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
