package plant

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/utils"
)

// PlantsHandler handles all requests to the plants endpoint.
func PlantsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HttpMethodNotAllowedResponse(w, "Allowed methods: GET")
		return
	}
	handlePlantsGet(w, r)
}

// PlantOverviewHandler handles all requests to the plant overview endpoint.
func PlantOverviewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HttpMethodNotAllowedResponse(w, "Allowed methods: GET")
		return
	}
	handlePlantOverviewGet(w, r)
}

// handlePlantsGet handles the retrieval of all plants.
func handlePlantsGet(w http.ResponseWriter, r *http.Request) {
	plantGroupIdStr := r.URL.Query().Get("plantGroupId")
	var filter *plantsFilter
	if plantGroupIdStr != "" {
		plantGroupId, err := strconv.ParseInt(plantGroupIdStr, 10, 64)
		if err != nil {
			msg := fmt.Sprintf("Error parsing plants filter: %s", err.Error())
			utils.HttpBadRequestResponse(w, msg)
			return
		}

		filter = &plantsFilter{
			PlantGroupId: plantGroupId,
		}
	}

	allPlants, err := getAllPlants(filter)
	if err != nil {
		msg := fmt.Sprintf("Error getting all plants: %s", err.Error())
		utils.HttpInternalServerErrorResponse(w, msg)
		return
	}

	b, err := json.Marshal(allPlants)
	if err != nil {
		msg := fmt.Sprintf("Error converting all plants to JSON: %s", err.Error())
		utils.HttpInternalServerErrorResponse(w, msg)
		return
	}

	log.Printf("Load %d plants", len(allPlants.Plants))
	utils.HttpOkResponse(w, b)
}

// getAllPlants retrieves all plants from the database.
func getAllPlants(filter *plantsFilter) (*plants, error) {
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

// handlePlantOverviewGet handles the retrieval of all plants.
func handlePlantOverviewGet(w http.ResponseWriter, r *http.Request) {
	allPlants, err := getAllPlantOverview()
	if err != nil {
		msg := fmt.Sprintf("Error getting all plants: %s", err.Error())
		utils.HttpInternalServerErrorResponse(w, msg)
		return
	}

	b, err := json.Marshal(allPlants)
	if err != nil {
		msg := fmt.Sprintf("Error converting all plants to JSON: %s", err.Error())
		utils.HttpInternalServerErrorResponse(w, msg)
		return
	}

	log.Printf("Load %d plants", len(allPlants.Plants))
	utils.HttpOkResponse(w, b)
}

// getAllPlantOverview retrieves all plants from the database.
func getAllPlantOverview() (*plantsOverview, error) {
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

	plantIds, err := plantRepository.GetAllOverview()
	return &plantsOverview{Plants: plantIds}, err
}
