package sensor_range

import "github.com/plantineers/plantbuddy-server/model"

type SensorRangeRepository interface {
	GetAllByPlantGroupId(id int64) ([]*model.SensorRange, error)
}
