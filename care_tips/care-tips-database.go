// Author: Maximilian Floto, Yannick Kirschen
package care_tips

// CareTipsRepository provides access to care tips.
type CareTipsRepository interface {
	// GetByPlantGroupId returns all care tips for a given plant group ID.
	GetByPlantGroupId(id int64) ([]string, error)

	// GetAdditionalByPlantId returns all additional care tips for a given plant ID.
	GetAdditionalByPlantId(id int64) ([]string, error)

	// Create creates new care tips for a given plant group ID.
	// Caution: This method does not use a transaction.
	Create(plantGroupId int64, careTips []string) error

	// DeleteAllByPlantGroupId deletes all care tips for a given plant group ID.
	// Caution: This method does not use a transaction.
	DeleteAllByPlantGroupId(id int64) error

	// CreateAdditionalByPlantId creates new additional care tips for a given plant ID.
	// Caution: This method does not use a transaction.
	CreateAdditionalByPlantId(plantId int64, careTips []string) error

	// DeleteAdditionalByPlantId deletes all additional care tips for a given plant ID.
	// Caution: This method does not use a transaction.
	DeleteAdditionalByPlantId(plantId int64) error
}
