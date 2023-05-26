// Author: Maximilian Floto, Yannick Kirschen
package plant

import "github.com/plantineers/plantbuddy-server/model"

type PlantRepository interface {
	// GetPlantById returns a plant by its ID.
	// If the plant does not exist, it will return nil.
	GetById(id int64) (*model.Plant, error)

	// Reads all plantIds from the database and returns them as a slice of plants.
	GetAll(filter *PlantsFilter) ([]int64, error)

	// Creates a new plant and returns its ID.
	Create(plant *model.PostPlantRequest) (*model.Plant, error)

	// Deletes a plant by its ID.
	DeleteById(id int64) error
}
