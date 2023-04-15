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

// Reads the buddy.json file and returns a Config object
func ReadConfig() (Config, error) {
	file, file_err := ioutil.ReadFile("buddy.json")
	if file_err != nil {
		return Config{}, file_err
	}
	var config Config
	json_err := json.Unmarshal(file, &config)
	if json_err != nil {
		return Config{}, json_err
	}
	return config, nil
}
