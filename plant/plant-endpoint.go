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

const convertPlantErrorStr = "Error converting plant %d to JSON: %s"

// PlantCreateHandler handles the creation of a new plant.
func PlantCreateHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlePlantPost(w, r)
	default:
		utils.HttpMethodNotAllowedResponse(w, "Allowed methods: POST")
	}
}

// PlantHandler handles all requests to the plant endpoint.
func PlantHandler(w http.ResponseWriter, r *http.Request) {
	id, err := utils.PathParameterFilter(r.URL.Path, "/v1/plant/")
	if err != nil {
		msg := fmt.Sprintf("Error getting path variable (plant ID): %s", err.Error())
		utils.HttpBadRequestResponse(w, msg)
		return
	}

	switch r.Method {
	case http.MethodGet:
		handlePlantGet(w, r, id)
	case http.MethodPut:
		handlePlantPut(w, r, id)
	case http.MethodDelete:
		handlePlantDelete(w, r, id)
	default:
		utils.HttpMethodNotAllowedResponse(w, "Allowed methods: GET, PUT, DELETE")
	}
}

// handlePlantPost handles the creation of a new plant.
func handlePlantPost(w http.ResponseWriter, r *http.Request) {
	var plant plantChange
	err := json.NewDecoder(r.Body).Decode(&plant)
	if err != nil {
		msg := fmt.Sprintf("Error decoding new plant: %s", err.Error())
		utils.HttpBadRequestResponse(w, msg)
		return
	}

	err = validate.Struct(plant)
	if err != nil {
		msg := fmt.Sprintf("Error validating new plant: %s", err.Error())
		utils.HttpBadRequestResponse(w, msg)
		return
	}

	createdPlantGroup, err := createPlant(&plant)
	if err == ErrPlantGroupNotExisting {
		msg := fmt.Sprintf("Plant group with id %d does not exist", plant.PlantGroupId)
		utils.HttpBadRequestResponse(w, msg)
		return
	}

	if err != nil {
		msg := fmt.Sprintf("Error creating plant: %s", err.Error())
		utils.HttpInternalServerErrorResponse(w, msg)
		return
	}

	b, err := json.Marshal(createdPlantGroup)
	if err != nil {
		msg := fmt.Sprintf(convertPlantErrorStr, createdPlantGroup.ID, err.Error())
		utils.HttpInternalServerErrorResponse(w, msg)
		return
	}

	msg := fmt.Sprintf("Plant with id %d created", createdPlantGroup.ID)
	location := fmt.Sprintf("/v1/plant/%d", createdPlantGroup.ID)
	utils.HttpCreatedResponse(w, b, location, msg)
}

// handlePlantGet handles the retrieval of a plant by its ID.
func handlePlantGet(w http.ResponseWriter, r *http.Request, id int64) {
	plant, err := getPlantById(id)
	switch err {
	case sql.ErrNoRows:
		msg := fmt.Sprintf("Plant with id %d not found", id)
		utils.HttpNotFoundResponse(w, msg)
	case nil:
		b, err := json.Marshal(plant)
		if err != nil {
			msg := fmt.Sprintf(convertPlantErrorStr, plant.ID, err.Error())
			utils.HttpInternalServerErrorResponse(w, msg)
			return
		}

		log.Printf("Plant with id %d loaded", id)
		utils.HttpOkResponse(w, b)
	default:
		msg := fmt.Sprintf("Error getting plant with id %d: %s", id, err.Error())
		utils.HttpBadRequestResponse(w, msg)
	}
}

// handlePlantPut handles the update of a plant by its ID.
func handlePlantPut(w http.ResponseWriter, r *http.Request, id int64) {
	var plantChange plantChange
	err := json.NewDecoder(r.Body).Decode(&plantChange)
	if err != nil {
		msg := fmt.Sprintf("Error decoding plant with id %d: %s", id, err.Error())
		utils.HttpBadRequestResponse(w, msg)
		return
	}

	err = validate.Struct(plantChange)
	if err != nil {
		msg := fmt.Sprintf("Error validating new plant: %s", err.Error())
		utils.HttpBadRequestResponse(w, msg)
		return
	}

	plant, err := updatePlantById(id, &plantChange)
	if err == ErrPlantGroupNotExisting {
		msg := fmt.Sprintf("Plant group with id %d does not exist", plantChange.PlantGroupId)
		utils.HttpBadRequestResponse(w, msg)
		return
	}

	if err != nil {
		msg := fmt.Sprintf("Error updating plant with id %d: %s", id, err.Error())
		utils.HttpInternalServerErrorResponse(w, msg)
		return
	}

	b, err := json.Marshal(plant)
	if err != nil {
		msg := fmt.Sprintf(convertPlantErrorStr, plant.ID, err.Error())
		utils.HttpInternalServerErrorResponse(w, msg)
		return
	}

	log.Printf("Plant with id %d updated", id)
	utils.HttpOkResponse(w, b)
}

// handlePlantDelete handles the deletion of a plant by its ID.
func handlePlantDelete(w http.ResponseWriter, r *http.Request, id int64) {
	err := deletePlantById(id)
	if err != nil {
		msg := fmt.Sprintf("Error deleting plant with id %d: %s", id, err.Error())
		utils.HttpInternalServerErrorResponse(w, msg)
		return
	}

	log.Printf("Plant with id %d deleted", id)
	utils.HttpOkResponse(w, nil)
}

// createPlant creates a new plant.
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

// getPlantById retrieves a plant by its ID.
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

// updatePlantById updates a plant by its ID.
func updatePlantById(id int64, plant *plantChange) (*Plant, error) {
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

	return repository.Update(id, plant)
}

// deletePlantById deletes a plant by its ID.
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
