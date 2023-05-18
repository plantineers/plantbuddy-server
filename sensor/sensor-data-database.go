// Author: Yannick Kirschen
package sensor

import "github.com/plantineers/plantbuddy-server/model"

type SensorDataRepository interface {
	GetAll(filter *SensorDataFilter) ([]*model.SensorData, error)

	Save(data *model.SensorData) error

	SaveAll(data []*model.SensorData) []error
}
