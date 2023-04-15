package config

import (
	"encoding/json"
	"os"
)

// Holds the configuration
type Config struct {
	Database Database
	Port     int `json:"port"`
}

// Holds the database configuration
type Database struct {
	DataSource string `json:"dataSource"`
	DriverName string `json:"driverName"`
}

// Holds the global configuration
var PlantBuddyConfig Config

// Reads the buddy.json file
func InitConfig() error {
	file, fileErr := os.ReadFile("buddy.json")
	if fileErr != nil {
		return fileErr
	}
	jsonErr := json.Unmarshal(file, &PlantBuddyConfig)
	if jsonErr != nil {
		return jsonErr
	}
	return nil
}
