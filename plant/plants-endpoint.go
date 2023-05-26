package plant

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/plantineers/plantbuddy-server/db"
)

// Returns all plants in database
func PlantsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	handlePlantsGet(w, r)
}

func handlePlantsGet(w http.ResponseWriter, r *http.Request) {
	plantGroupIdStr := r.URL.Query().Get("plantGroupId")
	var filter *PlantsFilter
	if plantGroupIdStr != "" {
		plantGroupId, err := strconv.ParseInt(plantGroupIdStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Error parsing plantGroupId: %s", err.Error())))
			return
		}

		filter = &PlantsFilter{
			PlantGroupId: plantGroupId,
		}
	}

	allPlants, err := getAllPlants(filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error getting all plants: %s", err.Error())))
		return
	}

	b, err := json.Marshal(allPlants)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error converting plants to JSON: %s", err.Error())))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func getAllPlants(filter *PlantsFilter) (*plants, error) {
	var session = db.NewSession()
	defer session.Close()

	err := session.Open()
	if err != nil {
		return nil, err
	}

	plantRepository, err := NewPlantRepository(session)
	if err != nil {
		return nil, err
	}

	plantIds, err := plantRepository.GetAll(filter)
	return &plants{Plants: plantIds}, err
}
