package plant

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/plantineers/plantbuddy-server/model"
	"github.com/plantineers/plantbuddy-server/utils"
)

func PlantHandler(w http.ResponseWriter, r *http.Request) {
	id, err := utils.PathParameterFilter(r.URL.Path, "/v1/plant/")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	switch r.Method {
	case http.MethodGet:
		plant, err := getPlantById(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		} else if plant == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Plant not found"))
			return
		}

		b, err := json.Marshal(plant)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error converting plant %d to JSON: %s", plant.ID, err.Error())))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
		return
	}
}

func getPlantById(id int64) (*model.Plant, error) {
	plantRepository := PlantSqliteRepository{}
	return plantRepository.GetById(id)
}
