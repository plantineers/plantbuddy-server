package sensor

import "github.com/plantineers/plantbuddy-server/model"

type SensorRepository interface {
	GetById(id int64) (*model.Sensor, error)

	GetAllIds() ([]int64, error)
}

type SensorTypeRepository interface {
	GetById(id int64) (*model.SensorType, error)

	GetAllIds() ([]int64, error)
}
