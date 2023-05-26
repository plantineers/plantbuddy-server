// Author: Yannick Kirschen
package sensor

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/plantineers/plantbuddy-server/db"
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

		log.Print(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
		return
	}

	b, err := json.Marshal(sensorTypes{Types: types})
	if err != nil {
		msg := fmt.Sprintf("Error converting sensor types to JSON: %s", err.Error())

		log.Print(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}

	log.Printf("Load %d sensor types", len(types))
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
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
