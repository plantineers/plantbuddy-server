// Author: Yannick Kirschen
package sensor

type SensorDataRepository interface {
	GetAll(filter *SensorDataFilter) ([]*SensorData, error)

	Save(data *SensorData) error

	SaveAll(data []*SensorData) []error
}
