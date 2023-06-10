package controller

// Controller represents a micro-controller, also called aggregator.
type Controller struct {
	UUID       string   `json:"uuid"`
	PlantGroup int64    `json:"plantGroup"`
	Sensors    []string `json:"sensors"`
}

// controllerUUIDs represents a list of controller UUIDs.
type controllerUUIDs struct {
	UUIDs []string `json:"controllers"`
}
