package sensor

import "github.com/plantineers/plantbuddy-server/model"

type SensorDataRepository interface {
	GetAll(filter *SensorDataFilter) ([]*model.SensorData, error)

	Save(data *model.SensorData) error
}

type SensorRepository interface {
	// GetById returns a sensor by its ID.
	// If the sensor does not exist, it will return nil.
	GetById(id int64) (*model.Sensor, error)

	// GetAllIds returns all sensor IDs.
	GetAllIds() ([]int64, error)
}

type SensorTypeRepository interface {
	// GetById returns a sensor type by its ID.
	// If the sensor does not exist, it will return nil.
	GetById(id int64) (*model.SensorType, error)

	// GetAllIds returns all sensor type IDs.
	GetAllIds() ([]int64, error)
}
