package plant

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/utils"
)

func PlantGroupHandler(w http.ResponseWriter, r *http.Request) {
	id, err := utils.PathParameterFilter(r.URL.Path, "/v1/plant-group/")
	if err != nil {
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

func PlantGroupCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	handlePlantGroupPost(w, r)
}

func handlePlantGroupGet(w http.ResponseWriter, r *http.Request, id int64) {
	plantGroup, err := getPlantGroupById(id)

	switch err {
	case sql.ErrNoRows:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Plant group not found"))
	case nil:
		b, err := json.Marshal(plantGroup)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error converting plant group %d to JSON: %s", plantGroup.ID, err.Error())))
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
}

func handlePlantGroupPut(w http.ResponseWriter, r *http.Request, id int64) {
	var plantGroup plantGroupChange
	err := json.NewDecoder(r.Body).Decode(&plantGroup)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error decoding JSON: %s", err.Error())))
		return
	}

	updatedPlantGroup, err := updatePlantGroup(id, &plantGroup)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error updating plant group: %s", err.Error())))
		return
	}

	b, err := json.Marshal(updatedPlantGroup)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error converting plant group %d to JSON: %s", updatedPlantGroup.ID, err.Error())))
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(b))
}

func handlePlantGroupDelete(w http.ResponseWriter, r *http.Request, id int64) {
	err := deletePlantGroup(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error deleting plant group: %s", err.Error())))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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

func handlePlantGroupPost(w http.ResponseWriter, r *http.Request) {
	var plantGroup plantGroupChange
	err := json.NewDecoder(r.Body).Decode(&plantGroup)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error decoding JSON: %s", err.Error())))
		return
	}

	createdPlantGroup, err := createPlantGroup(&plantGroup)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	b, err := json.Marshal(createdPlantGroup)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error converting plant group %d to JSON: %s", createdPlantGroup.ID, err.Error())))
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(b))
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
