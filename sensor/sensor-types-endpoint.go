// Author: Yannick Kirschen
package sensor

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/model"
)

func SensorTypesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleSensorTypesGet(w, r)
	}
}

func handleSensorTypesGet(w http.ResponseWriter, r *http.Request) {
	sensorTypes, err := getSensorTypes()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		b, err := json.Marshal(&model.SensorTypes{SensorTypes: sensorTypes})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error converting sensors to JSON: %s", err.Error())))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func getSensorTypes() ([]int64, error) {
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

	return repository.GetAllIds()
}
