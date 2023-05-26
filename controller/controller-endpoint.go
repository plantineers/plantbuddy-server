// Author: Yannick Kirschen
package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/utils"
)

func ControllerHandler(w http.ResponseWriter, r *http.Request) {
	uuid, err := utils.PathParameterFilterStr(r.URL.Path, "/v1/controller/")
	if err != nil {
		msg := fmt.Sprintf("Error getting path variable (controller UUID): %s", err.Error())

		log.Print(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
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
			msg := fmt.Sprintf("Error converting controller to JSON: %s", err.Error())

			log.Print(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(msg))
			return
		}

		log.Printf("Controller with UUID %s loaded", uuid)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	case sql.ErrNoRows:
		msg := fmt.Sprintf("Controller with UUID %s not found", uuid)

		log.Print(msg)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(msg))
	default:
		msg := fmt.Sprintf("Error getting controller with UUID %s: %s", uuid, err.Error())

		log.Print(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
	}
}

func getControllerData(uuid string) (*Controller, error) {
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
