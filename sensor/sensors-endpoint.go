// Author: Yannick Kirschen
package sensor

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/model"
)

func SensorsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleSensorsGet(w, r)
	}
}

func handleSensorsGet(w http.ResponseWriter, r *http.Request) {
	sensors, err := getSensors()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		b, err := json.Marshal(&model.Sensors{Sensors: sensors})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error converting sensors to JSON: %s", err.Error())))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func getSensors() ([]int64, error) {
	var session = db.NewSession()
	defer session.Close()

	err := session.Open()
	if err != nil {
		return nil, err
	}

	repository, err := NewSensorRepository(session)
	if err != nil {
		return nil, err
	}

	return repository.GetAllIds()
}
