// Author: Yannick Kirschen
package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/utils"
)

// ControllersHandler handles all requests to the controllers endpoint.
func ControllersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleControllersGet(w, r)
	}
}

// handleControllersGet handles GET requests to the controllers endpoint.
func handleControllersGet(w http.ResponseWriter, r *http.Request) {
	uuids, err := getAllControllerUUIDs()

	switch err {
	case nil:
		b, err := json.Marshal(&controllerUUIDs{UUIDs: uuids})
		if err != nil {
			msg := fmt.Sprintf("Error converting controller UUIDs to JSON: %s", err.Error())
			utils.HttpInternalServerErrorResponse(w, msg)
			return
		}

		utils.HttpOkResponse(w, b)
	default:
		msg := fmt.Sprintf("Error getting controller UUIDs: %s", err.Error())
		utils.HttpBadRequestResponse(w, msg)
	}
}

// getAllControllerUUIDs returns all UUIDs of all controllers.
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
