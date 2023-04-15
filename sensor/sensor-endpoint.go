package sensor

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/model"
	"github.com/plantineers/plantbuddy-server/utils"
)

func SensorHandler(w http.ResponseWriter, r *http.Request) {
	id, err := utils.PathParameterFilter(r.URL.Path, "/v1/sensor/")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleSensorGet(w, r, id)
	}
}

func SensorsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleSensorsGet(w, r)
	}
}

func handleSensorGet(w http.ResponseWriter, r *http.Request, id int64) {
	sensor, err := getSensorById(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else if sensor == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Sensor not found"))
	} else {
		b, err := json.Marshal(sensor)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error converting sensor %d to JSON: %s", sensor.ID, err.Error())))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
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

func getSensorById(id int64) (*model.Sensor, error) {
	var session = db.NewSession()
	defer session.Close()

	err := session.Open()
	if err != nil {
		return nil, err
	}

	repository, err := NewRepository(session)
	if err != nil {
		return nil, err
	}

	return repository.GetById(id)
}

func getSensors() ([]int64, error) {
	var session = db.NewSession()
	defer session.Close()

	err := session.Open()
	if err != nil {
		return nil, err
	}

	repository, err := NewRepository(session)
	if err != nil {
		return nil, err
	}

	return repository.GetAllIds()
}
