package plant

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/utils"
)

// PlantGroupsHandler handles all requests to the plant groups endpoint.
func PlantGroupsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HttpMethodNotAllowedResponse(w, "Allowed methods: GET")
		return
	}
	handlePlantGroupsGet(w, r)
}

// PlantGroupOverviewHandler handles all requests to the plant group overview endpoint.
func PlantGroupOverviewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HttpMethodNotAllowedResponse(w, "Allowed methods: GET")
		return
	}
	handlePlantGroupOverviewGet(w, r)
}

// handlePlantGroupsGet handles the retrieval of all plant groups.
func handlePlantGroupsGet(w http.ResponseWriter, r *http.Request) {
	allPlantGroups, err := getAllPlantGroups()
	if err != nil {
		msg := fmt.Sprintf("Error getting all plant groups: %s", err.Error())
		utils.HttpInternalServerErrorResponse(w, msg)
		return
	}

	b, err := json.Marshal(allPlantGroups)
	if err != nil {
		msg := fmt.Sprintf("Error converting all plant groups to JSON: %s", err.Error())
		utils.HttpInternalServerErrorResponse(w, msg)
		return
	}

	log.Printf("Load %d plant groups", len(allPlantGroups.PlantGroups))
	utils.HttpOkResponse(w, b)
}

// getAllPlantGroups retrieves all plant groups from the database.
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

// handlePlantGroupOverviewGet handles the retrieval of all plant group overviews.
func handlePlantGroupOverviewGet(w http.ResponseWriter, r *http.Request) {
	allPlantGroups, err := getAllPlantGroupOverview()
	if err != nil {
		msg := fmt.Sprintf("Error getting all plant groups: %s", err.Error())
		utils.HttpInternalServerErrorResponse(w, msg)
		return
	}

	b, err := json.Marshal(allPlantGroups)
	if err != nil {
		msg := fmt.Sprintf("Error converting all plant groups to JSON: %s", err.Error())
		utils.HttpInternalServerErrorResponse(w, msg)
		return
	}

	log.Printf("Load %d plant groups", len(allPlantGroups.PlantGroups))
	utils.HttpOkResponse(w, b)
}

// getAllPlantGroupOverview retrieves all plant group overviews from the database.
func getAllPlantGroupOverview() (*plantGroupsOverview, error) {
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

	plantGroupStubs, err := repository.GetAllOverview()
	return &plantGroupsOverview{PlantGroups: plantGroupStubs}, err
}
