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
		utils.HttpBadRequestResponse(w, msg)
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
			utils.HttpInternalServerErrorResponse(w, msg)
			return
		}

		log.Printf("Controller with UUID %s loaded", uuid)
		utils.HttpOkResponse(w, b)
	case sql.ErrNoRows:
		msg := fmt.Sprintf("Controller with UUID %s not found", uuid)
		utils.HttpNotFoundResponse(w, msg)
	default:
		msg := fmt.Sprintf("Error getting controller with UUID %s: %s", uuid, err.Error())
		utils.HttpBadRequestResponse(w, msg)
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
