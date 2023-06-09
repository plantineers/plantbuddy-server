// Author: Maximilian Floto, Yannick Kirschen
package plant

// PlantGroupRepository provides access to the plant group database.
type PlantGroupRepository interface {
	// GetPlantGroupById returns a plant group by its ID.
	GetById(id int64) (*PlantGroup, error)

	// Reads all plantGroupIds from the database and returns them as a slice of plant groups.
	GetAll() ([]int64, error)

	// Read all plant groups in the database and return them in short form (PlantGroupStub)
	GetAllOverview() ([]PlantGroupStub, error)

	// Create creates a new plant group in the database.
	Create(plantGroup *plantGroupChange) (*PlantGroup, error)

	// Update updates a plant group in the database.
	Update(id int64, plantGroup *plantGroupChange) (*PlantGroup, error)

	// Delete deletes a plant group from the database.
	Delete(id int64) error
}
