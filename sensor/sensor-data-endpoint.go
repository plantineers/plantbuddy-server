// Author: Yannick Kirschen
package sensor

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/plantineers/plantbuddy-server/db"
)

func SensorDataHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleSensorDataGet(w, r)
	case http.MethodPost:
		handleSensorDataPost(w, r)
	}
}

func handleSensorDataGet(w http.ResponseWriter, r *http.Request) {
	filter, err := filterSensorData(r)
	if err != nil {
		msg := fmt.Sprintf("Error parsing sensor data filter: %s", err.Error())

		log.Print(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
		return
	}

	allSensorData, err := getAllSensorData(filter)
	if err != nil {
		msg := fmt.Sprintf("Error getting all sensor data: %s", err.Error())

		log.Print(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}

	b, err := json.Marshal(&sensorDataSet{SensorData: allSensorData})
	if err != nil {
		msg := fmt.Sprintf("Error converting all sensor data to JSON: %s", err.Error())

		log.Print(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func handleSensorDataPost(w http.ResponseWriter, r *http.Request) {
	var data sensorDataPost

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		msg := fmt.Sprintf("Error parsing sensor data: %s", err.Error())

		log.Print(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}

	errs := saveSensorData(data.Data)
	switch errs {
	case nil:
		b, err := json.Marshal(data)
		if err != nil {
			msg := fmt.Sprintf("Error converting sensor data to JSON: %s", err.Error())

			log.Print(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(msg))
			return
		}

		log.Printf("Saved %d sensor data sets", len(data.Data))
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(b)
	default:
		var errStrings []string
		for _, err := range errs {
			errStrings = append(errStrings, err.Error())
		}

		msg := fmt.Sprintf("Error saving sensor data: %s", strings.Join(errStrings, "; "))

		log.Print(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
	}
}

func filterSensorData(r *http.Request) (*SensorDataFilter, error) {
	sensor := r.URL.Query().Get("sensor")
	plantStr := r.URL.Query().Get("plant")

	if sensor == "" || plantStr == "" {
		return nil, errors.New("sensor type and plant ID must be set")
	}

	plant, e := strconv.ParseInt(plantStr, 10, 64)
	if e != nil {
		return nil, errors.New("plant ID must be an integer")
	}

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	if from == "" {
		from = time.Now().AddDate(0, 0, -1).Format(time.RFC3339)
	}
	if to == "" {
		to = time.Now().Format(time.RFC3339)
	}

	return &SensorDataFilter{
		Sensor: sensor,
		Plant:  plant,
		From:   from,
		To:     to,
	}, nil
}

func getAllSensorData(filter *SensorDataFilter) ([]*SensorData, error) {
	var session = db.NewSession()
	defer session.Close()

	err := session.Open()
	if err != nil {
		return nil, err
	}

	repository, err := NewSensorDataRepository(session)
	if err != nil {
		return nil, err
	}

	return repository.GetAll(filter)
}

func saveSensorData(data []*SensorData) []error {
	var session = db.NewSession()
	defer session.Close()

	errors := make([]error, 0)
	err := session.Open()
	if err != nil {
		return append(errors, err)
	}

	repository, err := NewSensorDataRepository(session)
	if err != nil {
		return append(errors, err)
	}

	for _, d := range data {
		d.Timestamp = time.Now().UTC().String()
	}

	return repository.SaveAll(data)
}
