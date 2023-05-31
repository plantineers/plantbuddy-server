package sensor

func ValidateSensorRangeChange(sensorRange *SensorRangeChange) error {
	if sensorRange.Sensor == "" {
		return ErrSensorTypeRequired
	}
	return nil
}
