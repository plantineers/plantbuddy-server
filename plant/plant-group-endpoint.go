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

const convertPlantGroupErrorStr = "Error converting plant group %d to JSON: %s"

func PlantGroupCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HttpMethodNotAllowedResponse(w, "Allowed methods: POST")
		return
	}

	handlePlantGroupPost(w, r)
}

func PlantGroupHandler(w http.ResponseWriter, r *http.Request) {
	id, err := utils.PathParameterFilter(r.URL.Path, "/v1/plant-group/")
	if err != nil {
		msg := fmt.Sprintf("Error getting path variable (plant group ID): %s", err.Error())
		utils.HttpBadRequestResponse(w, msg)
		return
	}

	switch r.Method {
	case http.MethodGet:
		handlePlantGroupGet(w, r, id)
	case http.MethodPut:
		handlePlantGroupPut(w, r, id)
	case http.MethodDelete:
		handlePlantGroupDelete(w, r, id)
	default:
		utils.HttpMethodNotAllowedResponse(w, "Allowed methods: GET, PUT, DELETE")
	}
}

func handlePlantGroupPost(w http.ResponseWriter, r *http.Request) {
	var plantGroup plantGroupChange
	err := json.NewDecoder(r.Body).Decode(&plantGroup)
	if err != nil {
		msg := fmt.Sprintf("Error decoding new plant group: %s", err.Error())
		utils.HttpBadRequestResponse(w, msg)
		return
	}

	createdPlantGroup, err := createPlantGroup(&plantGroup)
	if err != nil {
		msg := fmt.Sprintf("Error creating new plant group: %s", err.Error())
		utils.HttpBadRequestResponse(w, msg)
		return
	}

	b, err := json.Marshal(createdPlantGroup)
	if err != nil {
		msg := fmt.Sprintf(convertPlantGroupErrorStr, createdPlantGroup.ID, err.Error())
		utils.HttpInternalServerErrorResponse(w, msg)
	}

	msg := fmt.Sprintf("Plant group with id %d created", createdPlantGroup.ID)
	location := fmt.Sprintf("/v1/plant-group/%d", createdPlantGroup.ID)
	utils.HttpCreatedResponse(w, b, location, msg)
}

func handlePlantGroupGet(w http.ResponseWriter, r *http.Request, id int64) {
	plantGroup, err := getPlantGroupById(id)

	switch err {
	case sql.ErrNoRows:
		msg := fmt.Sprintf("Plant group with id %d not found", id)
		utils.HttpNotFoundResponse(w, msg)
	case nil:
		b, err := json.Marshal(plantGroup)
		if err != nil {
			msg := fmt.Sprintf(convertPlantGroupErrorStr, plantGroup.ID, err.Error())
			utils.HttpInternalServerErrorResponse(w, msg)
		}

		log.Printf("Plant group with id %d loaded", id)
		utils.HttpOkResponse(w, b)
	default:
		msg := fmt.Sprintf("Error loading plant group with id %d: %s", id, err.Error())
		utils.HttpBadRequestResponse(w, msg)
	}
}

func handlePlantGroupPut(w http.ResponseWriter, r *http.Request, id int64) {
	var plantGroup plantGroupChange
	err := json.NewDecoder(r.Body).Decode(&plantGroup)
	if err != nil {
		msg := fmt.Sprintf("Error decoding plant group with id %d: %s", id, err.Error())
		utils.HttpBadRequestResponse(w, msg)
		return
	}

	updatedPlantGroup, err := updatePlantGroup(id, &plantGroup)
	if err != nil {
		msg := fmt.Sprintf("Error updating plant group with id %d: %s", id, err.Error())
		utils.HttpInternalServerErrorResponse(w, msg)
		return
	}

	b, err := json.Marshal(updatedPlantGroup)
	if err != nil {
		msg := fmt.Sprintf(convertPlantGroupErrorStr, updatedPlantGroup.ID, err.Error())
		utils.HttpInternalServerErrorResponse(w, msg)
	}

	log.Printf("Plant group with id %d updated", id)
	utils.HttpOkResponse(w, b)
}

func handlePlantGroupDelete(w http.ResponseWriter, r *http.Request, id int64) {
	err := deletePlantGroup(id)
	if err != nil {
		msg := fmt.Sprintf("Error deleting plant group with id %d: %s", id, err.Error())
		utils.HttpInternalServerErrorResponse(w, msg)
		return
	}

	log.Printf("Plant group with id %d deleted", id)
	utils.HttpOkResponse(w, nil)
}

func createPlantGroup(plantGroup *plantGroupChange) (*PlantGroup, error) {
	var session = db.NewSession()
	defer session.Close()

	err := session.Open()
	if err != nil {
		return nil, err
	}

	repository, err := NewPlantGroupRepository(session)
	if err != nil {
		return nil, err
	}

	return repository.Create(plantGroup)
}

func getPlantGroupById(id int64) (*PlantGroup, error) {
	var session = db.NewSession()
	defer session.Close()

	err := session.Open()
	if err != nil {
		return nil, err
	}

	repository, err := NewPlantGroupRepository(session)
	if err != nil {
		return nil, err
	}

	return repository.GetById(id)
}

func updatePlantGroup(id int64, plantGroup *plantGroupChange) (*PlantGroup, error) {
	var session = db.NewSession()
	defer session.Close()

	err := session.Open()
	if err != nil {
		return nil, err
	}

	repository, err := NewPlantGroupRepository(session)
	if err != nil {
		return nil, err
	}

	return repository.Update(id, plantGroup)
}

func deletePlantGroup(id int64) error {
	var session = db.NewSession()
	defer session.Close()

	err := session.Open()
	if err != nil {
		return err
	}

	repository, err := NewPlantGroupRepository(session)
	if err != nil {
		return err
	}

	return repository.Delete(id)
}
