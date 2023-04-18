// Package model contains the data model for the application`s REST interface.
//
// Author: Maximilian Floto, Yannick Kirschen
package model

type SensorData struct {
	ID        int64   `json:"id"`
	Sensor    int64   `json:"sensor"`
	Value     float64 `json:"value"`
	Timestamp string  `json:"timestamp"`
}

type Sensor struct {
	ID         int64       `json:"id"`
	Plant      int64       `json:"plant"`
	SensorType *SensorType `json:"sensorType"`
	Interval   int64       `json:"interval"`
}

type SensorRange struct {
	SensorType *SensorType `json:"sensorType"`
	Min        float64     `json:"min"`
	Max        float64     `json:"max"`
}

type SensorType struct {
	Name string `json:"name"`
	Unit string `json:"unit"`
}

type Sensors struct {
	Sensors []int64 `json:"sensors"`
}

type SensorTypes struct {
	SensorTypes []int64 `json:"sensorTypes"`
}

type Plant struct {
	ID                 int64       `json:"id"`
	Description        string      `json:"description"`
	PlantGroup         *PlantGroup `json:"plantGroup"`
	AdditionalCareTips []string    `json:"additionalCareTips"`
}

type Plants struct {
	Plants []int64 `json:"plants"`
}

type PlantGroup struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	CareTips    []string     `json:"careTips"`
	SensorRange *SensorRange `json:"sensorRange"`
}

type PlantGroups struct {
	PlantGroups []int64 `json:"plantGroups"`
}
