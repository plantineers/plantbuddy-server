package plant

import "github.com/plantineers/plantbuddy-server/sensor"

type Plant struct {
	ID                 int64       `json:"id"`
	Description        string      `json:"description"`
	Name               string      `json:"name"`
	Species            string      `json:"species"`
	Location           string      `json:"location"`
	PlantGroup         *PlantGroup `json:"plantGroup"`
	AdditionalCareTips []string    `json:"additionalCareTips"`
}

type PlantGroup struct {
	ID           int64                 `json:"id"`
	Name         string                `json:"name"`
	Description  string                `json:"description"`
	CareTips     []string              `json:"careTips"`
	SensorRanges []*sensor.SensorRange `json:"sensorRanges"`
}

type plants struct {
	Plants []int64 `json:"plants"`
}

type plantGroups struct {
	PlantGroups []int64 `json:"plantGroups"`
}

type plantChange struct {
	Description        string   `json:"description"`
	Name               string   `json:"name"`
	Species            string   `json:"species"`
	Location           string   `json:"location"`
	PlantGroupId       int64    `json:"plantGroupId"`
	AdditionalCareTips []string `json:"additionalCareTips"`
}

type plantGroupChange struct {
	Name         string                `json:"name"`
	Description  string                `json:"description"`
	CareTips     []string              `json:"careTips"`
	SensorRanges []*sensor.SensorRange `json:"sensorRanges"`
}
