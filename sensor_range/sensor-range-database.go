// Author: Yannick Kirschen
package sensor_range

import (
	"github.com/plantineers/plantbuddy-server/sensor"
)

type SensorRangeRepository interface {
	GetAllByPlantGroupId(id int64) ([]*sensor.SensorRange, error)

	Create(plantGroupId int64, sensorRange *sensor.SensorRange) error

	CreateAll(plantGroupId int64, sensorRanges []*sensor.SensorRange) error

	Update(plantGroupId int64, sensorRange *sensor.SensorRange) error

	UpdateAll(plantGroupId int64, sensorRanges []*sensor.SensorRange) error

	DeleteAllByPlantGroupId(id int64) error
}
