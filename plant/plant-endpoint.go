// Author: Maximilian Floto, Yannick Kirschen
package plant

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/plantineers/plantbuddy-server/db"
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
		handlePlantGet(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handlePlantGet(w http.ResponseWriter, r *http.Request, id int64) {
	plant, err := getPlantById(id)
	switch err {
	case sql.ErrNoRows:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Plant not found"))
	case nil:
		b, err := json.Marshal(plant)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error converting plant %d to JSON: %s", plant.ID, err.Error())))
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

	}
}

func getPlantById(id int64) (*model.Plant, error) {
	var session = db.NewSession()
	defer session.Close()

	err := session.Open()
	if err != nil {
		return nil, err
	}

	repository, err := NewPlantRepository(session)
	if err != nil {
		return nil, err
	}

	return repository.GetById(id)
}
