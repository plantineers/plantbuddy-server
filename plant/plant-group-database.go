// Author: Maximilian Floto, Yannick Kirschen
package plant

import "github.com/plantineers/plantbuddy-server/model"

type PlantGroupRepository interface {
	// GetPlantGroupById returns a plant group by its ID.
	GetById(id int64) (*model.PlantGroup, error)

	// Reads all plantGroupIds from the database and returns them as a slice of plant groups.
	GetAll() ([]int64, error)
}
