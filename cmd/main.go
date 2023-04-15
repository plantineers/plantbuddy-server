package main

import (
	"fmt"
	"net/http"

	"github.com/plantineers/plantbuddy-server/config"
	"github.com/plantineers/plantbuddy-server/plant"
)

func main() {

	// Read configuration file and panic if it fails
	config, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/v1/plant/", plant.PlantHandler)

	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}
