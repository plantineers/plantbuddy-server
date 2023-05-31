package plant

import "github.com/plantineers/plantbuddy-server/sensor"

func validatePlantChange(plant *plantChange) error {
	if plant.Name == "" {
		return ErrPlantNameRequired
	}
	if plant.PlantGroupId == 0 {
		return ErrPlantGroupRequired
	}
	return nil
}

func validatePlantGroupChange(plantGroup *plantGroupChange) error {
	if plantGroup.Name == "" {
		return ErrPlantGroupNameRequired
	}
	for _, sensorRange := range plantGroup.SensorRanges {
		err := sensor.ValidateSensorRangeChange(sensorRange)
		if err != nil {
			return err
		}
	}
	return nil
}
