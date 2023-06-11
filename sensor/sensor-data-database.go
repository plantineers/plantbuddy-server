// Author: Yannick Kirschen
package sensor

// SensorDataRepository provides access to sensor data.
type SensorDataRepository interface {
	// GetAll returns all sensor data matching the given filter.
	GetAll(filter *SensorDataFilter) ([]*SensorData, error)

	// Save stores the given sensor data.
	// Note: This is done using a transaction.
	Save(data *SensorData) error

	// SaveAll stores the given sensor data slice.
	// Note: Each data set is written in its own transaction.
	SaveAll(data []*SensorData) []error
}
