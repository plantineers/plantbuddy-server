// Author: Maximilian Floto, Yannick Kirschen
package plant

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/model"
	"github.com/plantineers/plantbuddy-server/utils"
)

func PlantHandler(w http.ResponseWriter, r *http.Request) {
	id, err := utils.PathParameterFilter(r.URL.Path, "/v1/plant/")
	if err != nil {
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

func PlantCreateHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlePlantPost(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handlePlantGet(w http.ResponseWriter, r *http.Request, id int64) {
	plant, err := getPlantById(id)
	switch err {
	case sql.ErrNoRows:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Plant not found"))
	case nil:
		b, err := json.Marshal(plant)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error converting plant %d to JSON: %s", plant.ID, err.Error())))
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

	}
}

func getPlantById(id int64) (*model.Plant, error) {
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

func handlePlantDelete(w http.ResponseWriter, r *http.Request, id int64) {
	err := deletePlantById(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error deleting plant: %s", err.Error())))
		return
	}

	w.WriteHeader(http.StatusOK)
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

func handlePlantPut(w http.ResponseWriter, r *http.Request, id int64) {
	var plant PlantChange
	err := json.NewDecoder(r.Body).Decode(&plant)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error decoding plant: %s", err.Error())))
		return
	}

	err = updatePlantById(id, &plant)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error updating plant: %s", err.Error())))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func updatePlantById(id int64, plant *PlantChange) error {
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

func handlePlantPost(w http.ResponseWriter, r *http.Request) {
	var plant PlantChange
	err := json.NewDecoder(r.Body).Decode(&plant)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error decoding plant: %s", err.Error())))
		return
	}

	newPlant, err := createPlant(&plant)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error creating plant: %s", err.Error())))
		return
	}

	b, err := json.Marshal(newPlant)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error converting plant %d to JSON: %s", newPlant.ID, err.Error())))
	}

	w.Header().Add("Content-Type", "application/json")

	w.Header().Add("Location", fmt.Sprintf("/v1/plant/%d", newPlant.ID))
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

func createPlant(plant *PlantChange) (*model.Plant, error) {
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

type PlantChange struct {
	Description        string   `json:"description"`
	Name               string   `json:"name"`
	Species            string   `json:"species"`
	Location           string   `json:"location"`
	PlantGroupId       int64    `json:"plantGroupId"`
	AdditionalCareTips []string `json:"additionalCareTips"`
}
