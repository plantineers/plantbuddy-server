package sensor

import (
	"database/sql"
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
	case http.MethodPut:
		handleSensorPut(w, r, id)
	case http.MethodDelete:
		handleSensorDelete(w, r, id)
	}
}

func handleSensorGet(w http.ResponseWriter, r *http.Request, id int64) {
	sensor, err := getSensorById(id)

	switch err {
	case sql.ErrNoRows:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Sensor not found"))
	case nil:
		b, err := json.Marshal(sensor)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error converting sensor %d to JSON: %s", sensor.ID, err.Error())))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
}

func getSensorById(id int64) (*model.Sensor, error) {
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

	return repository.GetById(id)
}

func handleSensorPut(w http.ResponseWriter, r *http.Request, id int64) {
	var sensor model.SensorPost
	err := json.NewDecoder(r.Body).Decode(&sensor)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error decoding JSON: %s", err.Error())))
		return
	}

	updatedSensor, err := updateSensor(&sensor, id)

	switch err {
	case sql.ErrNoRows:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Sensor not found"))
	case nil:
		b, err := json.Marshal(updatedSensor)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error converting sensor %d to JSON: %s", updatedSensor.ID, err.Error())))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

}

func updateSensor(sensor *model.SensorPost, id int64) (*model.Sensor, error) {
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

	return repository.Update(sensor, id)
}

func handleSensorDelete(w http.ResponseWriter, r *http.Request, id int64) {
	err := deleteSensor(id)

	switch err {
	case sql.ErrNoRows:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Sensor not found"))
	case nil:
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
}

func deleteSensor(id int64) error {
	var session = db.NewSession()
	defer session.Close()

	err := session.Open()
	if err != nil {
		return err
	}

	repository, err := NewSensorRepository(session)
	if err != nil {
		return err
	}

	return repository.Delete(id)
}

func SensorCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	handleSensorPost(w, r)
}

func handleSensorPost(w http.ResponseWriter, r *http.Request) {
	var sensor model.SensorPost
	err := json.NewDecoder(r.Body).Decode(&sensor)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error decoding JSON: %s", err.Error())))
		return
	}

	createdSensor, err := createSensor(&sensor)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	b, err := json.Marshal(createdSensor)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(b))
}

func createSensor(sensor *model.SensorPost) (*model.Sensor, error) {
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

	return repository.Create(sensor)
}
