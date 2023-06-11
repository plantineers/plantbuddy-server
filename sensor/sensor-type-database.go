// Author: Yannick Kirschen
package sensor

// SensorTypeRepository provides access to sensor types.
type SensorTypeRepository interface {
	// GetAll returns all sensor type IDs.
	GetAll() ([]*SensorType, error)
}
