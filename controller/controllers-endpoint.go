package controller

import (
	"encoding/json"
	"net/http"

	"github.com/plantineers/plantbuddy-server/db"
)

func ControllersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleControllersGet(w, r)
	}
}

func handleControllersGet(w http.ResponseWriter, r *http.Request) {
	uuids, err := getAllControllerUUIDs()

	switch err {
	case nil:
		b, err := json.Marshal(&controllerUUIDs{UUIDs: uuids})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
}

func getAllControllerUUIDs() ([]string, error) {
	var session = db.NewSession()
	defer session.Close()

	err := session.Open()
	if err != nil {
		return nil, err
	}

	repository, err := NewControllerRepository(session)
	if err != nil {
		return nil, err
	}

	return repository.GetAllUUIDs()
}

type controllerUUIDs struct {
	UUIDs []string `json:"controllers"`
}
