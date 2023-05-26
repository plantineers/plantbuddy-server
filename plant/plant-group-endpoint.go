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

func PlantGroupCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	handlePlantGroupPost(w, r)
}

func PlantGroupHandler(w http.ResponseWriter, r *http.Request) {
	id, err := utils.PathParameterFilter(r.URL.Path, "/v1/plant-group/")
	if err != nil {
		msg := fmt.Sprintf("Error getting path variable (plant group ID): %s", err.Error())

		log.Print(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
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
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handlePlantGroupPost(w http.ResponseWriter, r *http.Request) {
	var plantGroup plantGroupChange
	err := json.NewDecoder(r.Body).Decode(&plantGroup)
	if err != nil {
		msg := fmt.Sprintf("Error decoding new plant group: %s", err.Error())

		log.Print(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
		return
	}

	createdPlantGroup, err := createPlantGroup(&plantGroup)
	if err != nil {
		msg := fmt.Sprintf("Error creating new plant group: %s", err.Error())

		log.Print(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
		return
	}

	b, err := json.Marshal(createdPlantGroup)
	if err != nil {
		msg := fmt.Sprintf("Error converting plant group %d to JSON: %s", createdPlantGroup.ID, err.Error())

		log.Print(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
	}

	log.Printf("Plant group with id %d created", createdPlantGroup.ID)
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Location", fmt.Sprintf("/v1/plant-group/%d", createdPlantGroup.ID))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(b))
}

func handlePlantGroupGet(w http.ResponseWriter, r *http.Request, id int64) {
	plantGroup, err := getPlantGroupById(id)

	switch err {
	case sql.ErrNoRows:
		log.Printf("Plant group with id %d not found", id)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Plant group not found"))
	case nil:
		b, err := json.Marshal(plantGroup)
		if err != nil {
			msg := fmt.Sprintf("Error converting plant group %d to JSON: %s", plantGroup.ID, err.Error())

			log.Print(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(msg))
		}

		log.Printf("Plant group with id %d loaded", id)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	default:
		log.Printf("Error loading plant group with id %d: %s", id, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
}

func handlePlantGroupPut(w http.ResponseWriter, r *http.Request, id int64) {
	var plantGroup plantGroupChange
	err := json.NewDecoder(r.Body).Decode(&plantGroup)
	if err != nil {
		msg := fmt.Sprintf("Error decoding plant group with id %d: %s", id, err.Error())

		log.Print(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
		return
	}

	updatedPlantGroup, err := updatePlantGroup(id, &plantGroup)
	if err != nil {
		msg := fmt.Sprintf("Error updating plant group with id %d: %s", id, err.Error())

		log.Print(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}

	b, err := json.Marshal(updatedPlantGroup)
	if err != nil {
		msg := fmt.Sprintf("Error converting plant group %d to JSON: %s", updatedPlantGroup.ID, err.Error())

		log.Print(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
	}

	log.Printf("Plant group with id %d updated", id)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(b))
}

func handlePlantGroupDelete(w http.ResponseWriter, r *http.Request, id int64) {
	err := deletePlantGroup(id)
	if err != nil {
		msg := fmt.Sprintf("Error deleting plant group with id %d: %s", id, err.Error())

		log.Print(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}

	log.Printf("Plant group with id %d deleted", id)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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
