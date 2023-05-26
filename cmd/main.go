// The main executable of the PlantBuddy server.
//
// Author: Maximilian Floto, Yannick Kirschen
package main

import (
	"fmt"
	"net/http"

	"github.com/plantineers/plantbuddy-server/auth"
	"github.com/plantineers/plantbuddy-server/model"

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
	http.Handle("/v1/sensor-data", auth.UserAuthMiddleware(sensor.SensorDataHandler, model.Gardener, []string{"POST"}))

	http.Handle("/v1/sensor-types", auth.UserAuthMiddleware(sensor.SensorTypesHandler, model.Gardener, []string{}))

	http.Handle("/v1/controllers", auth.UserAuthMiddleware(controller.ControllersHandler, model.Gardener, []string{}))
	http.Handle("/v1/controller/", auth.UserAuthMiddleware(controller.ControllerHandler, model.Gardener, []string{}))

	http.Handle("/v1/plants", auth.UserAuthMiddleware(plant.PlantsHandler, model.Gardener, []string{}))
	http.Handle("/v1/plant", auth.UserAuthMiddleware(plant.PlantCreateHandler, model.Gardener, []string{}))
	http.Handle("/v1/plant/", auth.UserAuthMiddleware(plant.PlantHandler, model.Gardener, []string{}))

	http.Handle("/v1/plant-groups", auth.UserAuthMiddleware(plant.PlantGroupsHandler, model.Gardener, []string{}))
	http.Handle("/v1/plant-group", auth.UserAuthMiddleware(plant.PlantGroupCreateHandler, model.Gardener, []string{}))
	http.Handle("/v1/plant-group/", auth.UserAuthMiddleware(plant.PlantGroupHandler, model.Gardener, []string{}))

	http.Handle("/v1/users", auth.UserAuthMiddleware(auth.UsersHandler, model.Admin, []string{}))
	http.Handle("/v1/user/", auth.UserAuthMiddleware(auth.UserHandler, model.Admin, []string{}))
	http.HandleFunc("/v1/user/login", auth.LoginHandler)

	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", config.PlantBuddyConfig.Port), nil))
}
