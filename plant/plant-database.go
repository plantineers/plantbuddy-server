// Author: Maximilian Floto, Yannick Kirschen
package plant

type PlantRepository interface {
	// GetPlantById returns a plant by its ID.
	// If the plant does not exist, it will return nil.
	GetById(id int64) (*Plant, error)

	// Reads all plantIds from the database and returns them as a slice of plants.
	GetAll(filter *plantsFilter) ([]int64, error)

	// Creates a new plant and returns it.
	Create(plant *plantChange) (*Plant, error)

	// Updates a plant by its ID and returns it.
	Update(id int64, plant *plantChange) (*Plant, error)

	// Deletes a plant by its ID.
	DeleteById(id int64) error
}
