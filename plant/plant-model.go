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

type PlantStub struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type PlantGroup struct {
	ID           int64                 `json:"id"`
	Name         string                `json:"name"`
	Description  string                `json:"description"`
	CareTips     []string              `json:"careTips"`
	SensorRanges []*sensor.SensorRange `json:"sensorRanges"`
}

type PlantGroupStub struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type plants struct {
	Plants []int64 `json:"plants"`
}

type plantsOverview struct {
	Plants []PlantStub `json:"plants"`
}

type plantsFilter struct {
	PlantGroupId int64
}

type plantGroups struct {
	PlantGroups []int64 `json:"plantGroups"`
}

type plantGroupsOverview struct {
	PlantGroups []PlantGroupStub `json:"plantGroups"`
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
	Name         string                      `json:"name"`
	Description  string                      `json:"description"`
	CareTips     []string                    `json:"careTips"`
	SensorRanges []*sensor.SensorRangeChange `json:"sensorRanges"`
}
