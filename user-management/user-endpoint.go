package user_management

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/model"
	"github.com/plantineers/plantbuddy-server/utils"
	"net/http"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleUserGet(w, r)
	case http.MethodPost:
		handleUserPost(w, r)
	case http.MethodDelete:
		handleUserDelete(w, r)
	}
}

func handleUserGet(w http.ResponseWriter, r *http.Request) {
	name, err := utils.PathParameterFilterStr(r.URL.Path, "/v1/user/")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	user, err := getUserByName(name)
	switch err {
	case nil:
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

func handleUserPost(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.Password = utils.HashPassword(user.Password)

	err = createUser(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
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

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

func createUser(user *model.User) error {
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

	return repo.Create(user)
}

func handleUserDelete(w http.ResponseWriter, r *http.Request) {
	name, err := utils.PathParameterFilterStr(r.URL.Path, "/v1/user/")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = deleteUserByName(name)
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

func deleteUserByName(name string) error {
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

	user, err := getUserByName(name)
	if err != nil {
		return err
	}

	// Prevent admins from deleting other admins
	if user.Role == model.Admin {
		return ErrCannotDeleteAdmin
	}

	return repo.DeleteByName(name)
}

// TODO: Move to errors.go
var ErrCannotDeleteAdmin = errors.New("Cannot delete admin user")
