// Author: Yannick Kirschen
package sensor_range

import (
	"github.com/plantineers/plantbuddy-server/sensor"
)

type SensorRangeRepository interface {
	GetAllByPlantGroupId(id int64) ([]*sensor.SensorRange, error)

	Create(plantGroupId int64, sensorRange *sensor.SensorRangeChange) error

	CreateAll(plantGroupId int64, sensorRanges []*sensor.SensorRangeChange) error

	DeleteAllByPlantGroupId(id int64) error
}
