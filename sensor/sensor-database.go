package sensor

import "github.com/plantineers/plantbuddy-server/model"

type SensorRepository interface {
	GetById(id int64) (*model.Sensor, error)
}
