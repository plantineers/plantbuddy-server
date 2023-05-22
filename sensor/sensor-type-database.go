// Author: Yannick Kirschen
package sensor

import "github.com/plantineers/plantbuddy-server/model"

type SensorTypeRepository interface {
	// GetAll returns all sensor type IDs.
	GetAll() ([]*model.SensorType, error)
}
