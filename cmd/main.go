// The main executable of the PlantBuddy server.
//
// Author: Maximilian Floto, Yannick Kirschen
package main

import (
	"fmt"
	"github.com/plantineers/plantbuddy-server/model"
	user_management "github.com/plantineers/plantbuddy-server/user-management"
	"net/http"

	"github.com/plantineers/plantbuddy-server/config"
	"github.com/plantineers/plantbuddy-server/controller"
	"github.com/plantineers/plantbuddy-server/plant"
	"github.com/plantineers/plantbuddy-server/sensor"
)

func main() {
	// Read configuration file into a global variable and panic if it fails
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/v1/sensor-data", sensor.SensorDataHandler)

	http.HandleFunc("/v1/sensor-types", sensor.SensorTypesHandler)

	http.HandleFunc("/v1/controllers", controller.ControllersHandler)
	http.HandleFunc("/v1/controller/", controller.ControllerHandler)

	http.HandleFunc("/v1/plants", plant.PlantsHandler)
	http.HandleFunc("/v1/plant/", plant.PlantHandler)

	http.HandleFunc("/v1/plant-groups", plant.PlantGroupsHandler)
	http.HandleFunc("/v1/plant-group/", plant.PlantGroupHandler)

	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", config.PlantBuddyConfig.Port), nil))
}
