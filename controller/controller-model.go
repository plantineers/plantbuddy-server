package controller

type Controller struct {
	UUID       string   `json:"uuid"`
	PlantGroup int64    `json:"plantGroup"`
	Sensors    []string `json:"sensors"`
}
