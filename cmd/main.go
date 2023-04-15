package main

import (
	"fmt"
	"net/http"

	"github.com/plantineers/plantbuddy-server/plant"
)

func main() {
	http.HandleFunc("/v1/plant/", plant.PlantHandler)

	fmt.Println(http.ListenAndServe(":3333", nil))
}
