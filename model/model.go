// Package model contains the data model for the application`s REST interface.
//
// Author: Maximilian Floto, Yannick Kirschen
package model

type SensorData struct {
	Controller string  `json:"controller"`
	Sensor     string  `json:"sensor"`
	Value      float64 `json:"value"`
	Timestamp  string  `json:"timestamp"`
}

type Controller struct {
	UUID    string   `json:"uuid"`
	Plant   int64    `json:"plant"`
	Sensors []string `json:"sensors"`
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

type SensorTypes struct {
	Types []*SensorType `json:"types"`
}

type Plant struct {
	ID                 int64       `json:"id"`
	Description        string      `json:"description"`
	Name               string      `json:"name"`
	Species            string      `json:"species"`
	Location           string      `json:"location"`
	PlantGroup         *PlantGroup `json:"plantGroup"`
	AdditionalCareTips []string    `json:"additionalCareTips"`
}

type Plants struct {
	Plants []int64 `json:"plants"`
}

type PostPlantRequest struct {
	Description        string   `json:"description"`
	Name               string   `json:"name"`
	Species            string   `json:"species"`
	Location           string   `json:"location"`
	PlantGroupId       int64    `json:"plantGroupId"`
	AdditionalCareTips []string `json:"additionalCareTips"`
}

type PlantGroup struct {
	ID           int64          `json:"id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	CareTips     []string       `json:"careTips"`
	SensorRanges []*SensorRange `json:"sensorRanges"`
}

type PlantGroups struct {
	PlantGroups []int64 `json:"plantGroups"`
}

type User struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
}

type SafeUser struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Role Role   `json:"role"`
}

type Users struct {
	Users []string `json:"users"`
}

type Role int8

const (
	Admin Role = iota
	Gardener
)
