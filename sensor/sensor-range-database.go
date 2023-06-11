// Author: Yannick Kirschen
package sensor

// SensorRangeRepository provides access to sensor ranges.
type SensorRangeRepository interface {
	// GetAllByPlantGroupId returns all sensor ranges for the given plant group.
	GetAllByPlantGroupId(id int64) ([]*SensorRange, error)

	// Create stores the given sensor range and associates it with the given plant group.
	// Caution: This method does not use a transaction.
	Create(plantGroupId int64, sensorRange *SensorRangeChange) error

	// CreateAll stores the given sensor ranges and associates them with the given plant group.
	// Note: This method uses a transaction.
	CreateAll(plantGroupId int64, sensorRanges []*SensorRangeChange) error

	// DeleteAllByPlantGroupId deletes all sensor ranges associated with the given plant group.
	// Caution: This method does not use a transaction.
	DeleteAllByPlantGroupId(id int64) error
}
