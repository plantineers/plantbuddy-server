package plant

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/model"
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
	allPlants, err := getAllPlants()
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
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func getAllPlants() (*[]*model.Plant, error) {
	var session = db.NewSession()
	defer session.Close()

	err := session.Open()
	if err != nil {
		return nil, err
	}

	plantRepository, err := NewRepository(session)
	if err != nil {
		return nil, err
	}

	return plantRepository.GetAllPlants()
}
