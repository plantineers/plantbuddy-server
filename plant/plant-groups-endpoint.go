package plant

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/plantineers/plantbuddy-server/db"
)

func PlantGroupsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	handlePlantGroupsGet(w, r)
}

func handlePlantGroupsGet(w http.ResponseWriter, r *http.Request) {
	allPlantGroups, err := getAllPlantGroups()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error getting all plant groups: %s", err.Error())))
		return
	}
	b, err := json.Marshal(allPlantGroups)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error converting plant groups to JSON: %s", err.Error())))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func getAllPlantGroups() (*plantGroups, error) {
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

	plantGroupIds, err := repository.GetAll()
	return &plantGroups{PlantGroups: plantGroupIds}, err
}
