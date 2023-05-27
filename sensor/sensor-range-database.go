// Author: Yannick Kirschen
package sensor

type SensorRangeRepository interface {
	GetAllByPlantGroupId(id int64) ([]*SensorRange, error)

	Create(plantGroupId int64, sensorRange *SensorRangeChange) error

	CreateAll(plantGroupId int64, sensorRanges []*SensorRangeChange) error

	DeleteAllByPlantGroupId(id int64) error
}
