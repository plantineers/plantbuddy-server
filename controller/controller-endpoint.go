// Author: Yannick Kirschen
package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/model"
	"github.com/plantineers/plantbuddy-server/utils"
)

func ControllerHandler(w http.ResponseWriter, r *http.Request) {
	uuid, err := utils.PathParameterFilterStr(r.URL.Path, "/v1/controller/")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleControllerGet(w, r, uuid)
	}
}

func handleControllerGet(w http.ResponseWriter, r *http.Request, uuid string) {
	controller, err := getControllerData(uuid)
	switch err {
	case nil:
		b, err := json.Marshal(controller)
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
		w.Write([]byte("Controller not found"))
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
}

func getControllerData(uuid string) (*model.Controller, error) {
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

	return repository.GetByUUID(uuid)
}
