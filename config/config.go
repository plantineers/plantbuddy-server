package config

import (
	"encoding/json"
	"io/ioutil"
)

// Holds the configuration
type Config struct {
	Database Database
	Port     int `json:"port"`
}

// Holds the database configuration
type Database struct {
	DataSource string `json:"data_source"`
	DriverName string `json:"driver_name"`
}

// Holds the global configuration
var PlantBuddyConfig Config

// Reads the buddy.json file
func InitConfig() error {
	file, file_err := ioutil.ReadFile("buddy.json")
	if file_err != nil {
		return file_err
	}
	json_err := json.Unmarshal(file, &PlantBuddyConfig)
	if json_err != nil {
		return json_err
	}
	return nil
}
