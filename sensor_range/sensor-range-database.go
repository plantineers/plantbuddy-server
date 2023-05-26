package sensor_range

import (
	"github.com/plantineers/plantbuddy-server/model"
)

type SensorRangeRepository interface {
	GetAllByPlantGroupId(id int64) ([]*model.SensorRange, error)

	Create(plantGroupId int64, sensorRange *model.SensorRange) error

	CreateAll(plantGroupId int64, sensorRanges []*model.SensorRange) error

	Update(plantGroupId int64, sensorRange *model.SensorRange) error

	UpdateAll(plantGroupId int64, sensorRanges []*model.SensorRange) error

	DeleteAllByPlantGroupId(id int64) error
}
