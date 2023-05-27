// Author: Yannick Kirschen
package sensor

type SensorRangeRepository interface {
	GetAllByPlantGroupId(id int64) ([]*SensorRange, error)

	Create(plantGroupId int64, sensorRange *SensorRange) error

	CreateAll(plantGroupId int64, sensorRanges []*SensorRange) error

	Update(plantGroupId int64, sensorRange *SensorRange) error

	UpdateAll(plantGroupId int64, sensorRanges []*SensorRange) error

	DeleteAllByPlantGroupId(id int64) error
}
