package model

type Sensor struct {
	ID         int64       `json:"id"`
	Plant      *Plant      `json:"plant"`
	SensorType *SensorType `json:"sensorType"`
	Interval   int64       `json:"interval"`
}

type SensorRange struct {
	SensorType *SensorType `json:"sensorType"`
	Min        float64     `json:"min"`
	Max        float64     `json:"max"`
}

type Plant struct {
	ID                 int64       `json:"id"`
	Description        string      `json:"description"`
	PlantGroup         *PlantGroup `json:"plantGroup"`
	AdditionalCareTips []string    `json:"additionalCareTips"`
}

type PlantGroup struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	CareTips    []string     `json:"careTips"`
	SensorRange *SensorRange `json:"sensorRange"`
}

type SensorType struct {
	Name string `json:"name"`
	Unit string `json:"unit"`
}
