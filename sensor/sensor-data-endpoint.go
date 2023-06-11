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
	"github.com/plantineers/plantbuddy-server/utils"
)

// SensorDataHandler handles requests to the sensor-data endpoint.
func SensorDataHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleSensorDataGet(w, r)
	case http.MethodPost:
		handleSensorDataPost(w, r)
	}
}

// handleSensorDataGet handles GET requests to the sensor-data endpoint.
func handleSensorDataGet(w http.ResponseWriter, r *http.Request) {
	filter, err := filterSensorData(r)
	if err != nil {
		msg := fmt.Sprintf("Error parsing sensor data filter: %s", err.Error())
		utils.HttpBadRequestResponse(w, msg)
		return
	}

	allSensorData, err := getAllSensorData(filter)
	if err != nil {
		msg := fmt.Sprintf("Error getting all sensor data: %s", err.Error())
		utils.HttpInternalServerErrorResponse(w, msg)
		return
	}

	b, err := json.Marshal(&sensorDataSet{SensorData: allSensorData})
	if err != nil {
		msg := fmt.Sprintf("Error converting all sensor data to JSON: %s", err.Error())
		utils.HttpInternalServerErrorResponse(w, msg)
		return
	}

	utils.HttpOkResponse(w, b)
}

// handleSensorDataPost handles POST requests to the sensor-data endpoint.
func handleSensorDataPost(w http.ResponseWriter, r *http.Request) {
	var data sensorDataPost

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		msg := fmt.Sprintf("Error parsing sensor data: %s", err.Error())
		utils.HttpInternalServerErrorResponse(w, msg)
		return
	}

	errs := saveSensorData(data.Data)
	switch errs {
	case nil:
		b, err := json.Marshal(data)
		if err != nil {
			msg := fmt.Sprintf("Error converting sensor data to JSON: %s", err.Error())
			utils.HttpInternalServerErrorResponse(w, msg)
			return
		}

		log.Printf("Saved %d sensor data sets", len(data.Data))
		utils.HttpOkResponse(w, b)
	default:
		var errStrings []string
		for _, err := range errs {
			errStrings = append(errStrings, err.Error())
		}

		msg := fmt.Sprintf("Error saving sensor data: %s", strings.Join(errStrings, "; "))
		utils.HttpBadRequestResponse(w, msg)
	}
}

// filterSensorData parses the query parameters of a request and returns a SensorDataFilter.
func filterSensorData(r *http.Request) (*SensorDataFilter, error) {
	sensor := r.URL.Query().Get("sensor")
	plantStr := r.URL.Query().Get("plant")
	plantGroupStr := r.URL.Query().Get("plantGroup")

	if sensor == "" || (plantStr == "" && plantGroupStr == "") {
		return nil, errors.New("sensor type and either plant ID or plantGroup ID must be set")
	}

	var plant int64
	var plantGroup int64
	var err error

	if plantStr != "" {
		plant, err = strconv.ParseInt(plantStr, 10, 64)
		if err != nil {
			return nil, errors.New("plant ID must be an integer")
		}
	}

	if plantGroupStr != "" {
		plantGroup, err = strconv.ParseInt(plantGroupStr, 10, 64)
		if err != nil {
			return nil, errors.New("plantGroup ID must be an integer")
		}
	}

	if plant != 0 && plantGroup != 0 {
		return nil, errors.New("plant ID and plantGroup ID cannot be set at the same time")
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
		Sensor:     sensor,
		Plant:      plant,
		PlantGroup: plantGroup,
		From:       from,
		To:         to,
	}, nil
}

// getAllSensorData returns all sensor data matching the given filter.
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

// saveSensorData saves the given sensor data.
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
