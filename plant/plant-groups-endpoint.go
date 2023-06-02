package plant

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/utils"
)

func PlantGroupsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HttpMethodNotAllowedResponse(w, "Allowed methods: GET")
		return
	}
	handlePlantGroupsGet(w, r)
}

func PlantGroupOverviewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HttpMethodNotAllowedResponse(w, "Allowed methods: GET")
		return
	}
	handlePlantGroupOverviewGet(w, r)
}

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

	log.Printf("Load %d plant groups", len(b))
	utils.HttpOkResponse(w, b)
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

	log.Printf("Load %d plant groups", len(b))
	utils.HttpOkResponse(w, b)
}

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
