// Author: Yannick Kirschen
package sensor

type SensorTypeRepository interface {
	// GetAll returns all sensor type IDs.
	GetAll() ([]*SensorType, error)
}
