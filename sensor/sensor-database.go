package sensor

import "github.com/plantineers/plantbuddy-server/model"

type SensorRepository interface {
	// GetById returns a sensor by its ID.
	// If the sensor does not exist, it will return nil.
	GetById(id int64) (*model.Sensor, error)

	// GetAllIds returns all sensor IDs.
	GetAllIds() ([]int64, error)

	// Create creates a new sensor and returns it.
	Create(*model.SensorPost) (*model.Sensor, error)

	// Update updates an existing sensor and returns it.
	Update(*model.SensorPost, int64) (*model.Sensor, error)
}

type SensorTypeRepository interface {
	// GetById returns a sensor type by its ID.
	// If the sensor does not exist, it will return nil.
	GetById(id int64) (*model.SensorType, error)

	// GetAllIds returns all sensor type IDs.
	GetAllIds() ([]int64, error)
}
