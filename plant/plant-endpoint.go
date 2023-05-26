// Author: Maximilian Floto, Yannick Kirschen
package plant

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/utils"
)

func PlantCreateHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlePlantPost(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func PlantHandler(w http.ResponseWriter, r *http.Request) {
	id, err := utils.PathParameterFilter(r.URL.Path, "/v1/plant/")
	if err != nil {
		msg := fmt.Sprintf("Error getting path variable (plant ID): %s", err.Error())

		log.Print(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	switch r.Method {
	case http.MethodGet:
		handlePlantGet(w, r, id)
	case http.MethodDelete:
		handlePlantDelete(w, r, id)
	case http.MethodPut:
		handlePlantPut(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handlePlantPost(w http.ResponseWriter, r *http.Request) {
	var plant plantChange
	err := json.NewDecoder(r.Body).Decode(&plant)
	if err != nil {
		msg := fmt.Sprintf("Error decoding new plant: %s", err.Error())

		log.Print(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
		return
	}

	createdPlantGroup, err := createPlant(&plant)
	if err != nil {
		msg := fmt.Sprintf("Error creating plant: %s", err.Error())

		log.Print(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}

	b, err := json.Marshal(createdPlantGroup)
	if err != nil {
		msg := fmt.Sprintf("Error converting plant %d to JSON: %s", createdPlantGroup.ID, err.Error())

		log.Print(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
	}

	log.Printf("Plant with id %d created", createdPlantGroup.ID)
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Location", fmt.Sprintf("/v1/plant/%d", createdPlantGroup.ID))
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

func handlePlantGet(w http.ResponseWriter, r *http.Request, id int64) {
	plant, err := getPlantById(id)
	switch err {
	case sql.ErrNoRows:
		log.Printf("Plant with id %d not found", id)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Plant not found"))
	case nil:
		b, err := json.Marshal(plant)
		if err != nil {
			msg := fmt.Sprintf("Error converting plant %d to JSON: %s", plant.ID, err.Error())

			log.Print(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(msg))
		}

		log.Printf("Plant with id %d loaded", id)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	default:
		log.Printf("Error getting plant with id %d: %s", id, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
}

func handlePlantPut(w http.ResponseWriter, r *http.Request, id int64) {
	var plant plantChange
	err := json.NewDecoder(r.Body).Decode(&plant)
	if err != nil {
		msg := fmt.Sprintf("Error decoding plant with id %d: %s", id, err.Error())

		log.Print(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
		return
	}

	err = updatePlantById(id, &plant)
	if err != nil {
		msg := fmt.Sprintf("Error updating plant with id %d: %s", id, err.Error())

		log.Print(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}

	log.Printf("Plant with id %d updated", id)
	w.WriteHeader(http.StatusOK)
}

func handlePlantDelete(w http.ResponseWriter, r *http.Request, id int64) {
	err := deletePlantById(id)
	if err != nil {
		msg := fmt.Sprintf("Error deleting plant with id %d: %s", id, err.Error())

		log.Print(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}

	log.Printf("Plant with id %d deleted", id)
	w.WriteHeader(http.StatusOK)
}

func createPlant(plant *plantChange) (*Plant, error) {
	var session = db.NewSession()
	defer session.Close()

	err := session.Open()
	if err != nil {
		return nil, err
	}

	repository, err := NewPlantRepository(session)
	if err != nil {
		return nil, err
	}

	return repository.Create(plant)
}

func getPlantById(id int64) (*Plant, error) {
	var session = db.NewSession()
	defer session.Close()

	err := session.Open()
	if err != nil {
		return nil, err
	}

	repository, err := NewPlantRepository(session)
	if err != nil {
		return nil, err
	}

	return repository.GetById(id)
}

func updatePlantById(id int64, plant *plantChange) error {
	var session = db.NewSession()
	defer session.Close()

	err := session.Open()
	if err != nil {
		return err
	}

	repository, err := NewPlantRepository(session)
	if err != nil {
		return err
	}

	return repository.Update(id, plant)
}

func deletePlantById(id int64) error {
	var session = db.NewSession()
	defer session.Close()

	err := session.Open()
	if err != nil {
		return err
	}

	repository, err := NewPlantRepository(session)
	if err != nil {
		return err
	}

	return repository.DeleteById(id)
}
