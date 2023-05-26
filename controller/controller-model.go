package controller

type Controller struct {
	UUID    string   `json:"uuid"`
	Plant   int64    `json:"plant"`
	Sensors []string `json:"sensors"`
}
