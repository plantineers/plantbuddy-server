// Author: Yannick Kirschen
package sensor

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/utils"
)

func SensorTypesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleSensorTypesGet(w, r)
	}
}

func handleSensorTypesGet(w http.ResponseWriter, r *http.Request) {
	types, err := getSensorTypes()
	if err != nil {
		msg := fmt.Sprintf("Error getting sensor types: %s", err.Error())
		utils.HttpBadRequestResponse(w, msg)
		return
	}

	b, err := json.Marshal(sensorTypes{Types: types})
	if err != nil {
		msg := fmt.Sprintf("Error converting sensor types to JSON: %s", err.Error())
		utils.HttpInternalServerErrorResponse(w, msg)
		return
	}

	log.Printf("Load %d sensor types", len(types))
	utils.HttpOkResponse(w, b)
}

func getSensorTypes() ([]*SensorType, error) {
	var session = db.NewSession()
	defer session.Close()

	err := session.Open()
	if err != nil {
		return nil, err
	}

	repository, err := NewSensorTypeRepository(session)
	if err != nil {
		return nil, err
	}

	return repository.GetAll()
}
