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

	http.Handle("/v1/sensors", user_management.UserAuthMiddleware(sensor.SensorsHandler, model.Gardener))
	http.Handle("/v1/sensor", user_management.UserAuthMiddleware(sensor.SensorCreateHandler, model.Admin))
	http.HandleFunc("/v1/sensor/", sensor.SensorHandler)

	http.HandleFunc("/v1/sensor-types", sensor.SensorTypesHandler)
	http.HandleFunc("/v1/sensor-type/", sensor.SensorTypeHandler)

	http.Handle("/v1/plants", user_management.UserAuthMiddleware(plant.PlantsHandler, model.Gardener))
	http.Handle("/v1/plant/", user_management.UserAuthMiddleware(plant.PlantHandler, model.Gardener))

	http.Handle("/v1/plant-groups", user_management.UserAuthMiddleware(plant.PlantGroupsHandler, model.Gardener))
	http.Handle("/v1/plant-group/", user_management.UserAuthMiddleware(plant.PlantGroupHandler, model.Gardener))

	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", config.PlantBuddyConfig.Port), nil))
}
