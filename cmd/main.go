// The main executable of the PlantBuddy server.
//
// Author: Maximilian Floto, Yannick Kirschen
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/plantineers/plantbuddy-server/auth"
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

	// Initialize the validator for the plant package
	plant.InitializeValidator()

	http.Handle("/v1/sensor-data", auth.UserAuthMiddleware(sensor.SensorDataHandler, auth.Gardener))

	http.Handle("/v1/sensor-types", auth.UserAuthMiddleware(sensor.SensorTypesHandler, auth.Gardener))

	http.Handle("/v1/controllers", auth.UserAuthMiddleware(controller.ControllersHandler, auth.Gardener))
	http.Handle("/v1/controller/", auth.UserAuthMiddleware(controller.ControllerHandler, auth.Gardener))

	http.Handle("/v1/plants", auth.UserAuthMiddleware(plant.PlantsHandler, auth.Gardener))
	http.Handle("/v1/plants/overview", auth.UserAuthMiddleware(plant.PlantOverviewHandler, auth.Gardener))
	http.Handle("/v1/plant", auth.UserAuthMiddleware(plant.PlantCreateHandler, auth.Gardener))
	http.Handle("/v1/plant/", auth.UserAuthMiddleware(plant.PlantHandler, auth.Gardener))

	http.Handle("/v1/plant-groups", auth.UserAuthMiddleware(plant.PlantGroupsHandler, auth.Gardener))
	http.Handle("/v1/plant-groups/overview", auth.UserAuthMiddleware(plant.PlantGroupOverviewHandler, auth.Gardener))
	http.Handle("/v1/plant-group", auth.UserAuthMiddleware(plant.PlantGroupCreateHandler, auth.Gardener))
	http.Handle("/v1/plant-group/", auth.UserAuthMiddleware(plant.PlantGroupHandler, auth.Gardener))

	http.Handle("/v1/users", auth.UserAuthMiddleware(auth.UsersHandler, auth.Admin))
	http.Handle("/v1/user", auth.UserAuthMiddleware(auth.UserCreateHandler, auth.Admin))
	http.Handle("/v1/user/", auth.UserAuthMiddleware(auth.UserHandler, auth.Admin))
	http.HandleFunc("/v1/user/login", auth.LoginHandler)

	log.Printf("Server running on port %d", config.PlantBuddyConfig.Port)
	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", config.PlantBuddyConfig.Port), nil))
}
