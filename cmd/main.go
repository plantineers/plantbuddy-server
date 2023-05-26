// The main executable of the PlantBuddy server.
//
// Author: Maximilian Floto, Yannick Kirschen
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/plantineers/plantbuddy-server/model"
	user_management "github.com/plantineers/plantbuddy-server/user-management"

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

	// The POST method is not subject to user authentication as it is used by aggregators to send data
	http.Handle("/v1/sensor-data", user_management.UserAuthMiddleware(sensor.SensorDataHandler, model.Gardener, []string{"POST"}))

	http.Handle("/v1/sensor-types", user_management.UserAuthMiddleware(sensor.SensorTypesHandler, model.Gardener, []string{}))

	http.Handle("/v1/controllers", user_management.UserAuthMiddleware(controller.ControllersHandler, model.Gardener, []string{}))
	http.Handle("/v1/controller/", user_management.UserAuthMiddleware(controller.ControllerHandler, model.Gardener, []string{}))

	http.Handle("/v1/plants", user_management.UserAuthMiddleware(plant.PlantsHandler, model.Gardener, []string{}))
	http.Handle("/v1/plant", user_management.UserAuthMiddleware(plant.PlantCreateHandler, model.Gardener, []string{}))
	http.Handle("/v1/plant/", user_management.UserAuthMiddleware(plant.PlantHandler, model.Gardener, []string{}))

	http.Handle("/v1/plant-groups", user_management.UserAuthMiddleware(plant.PlantGroupsHandler, model.Gardener, []string{}))
	http.Handle("/v1/plant-group", user_management.UserAuthMiddleware(plant.PlantGroupCreateHandler, model.Gardener, []string{}))
	http.Handle("/v1/plant-group/", user_management.UserAuthMiddleware(plant.PlantGroupHandler, model.Gardener, []string{}))

	http.Handle("/v1/users", user_management.UserAuthMiddleware(user_management.UsersHandler, model.Admin, []string{}))
	http.Handle("/v1/user/", user_management.UserAuthMiddleware(user_management.UserHandler, model.Admin, []string{}))
	http.HandleFunc("/v1/user/login", user_management.LoginHandler)

	log.Printf("Server running on port %d", config.PlantBuddyConfig.Port)
	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", config.PlantBuddyConfig.Port), nil))
}
