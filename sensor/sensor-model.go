package sensor

type SensorData struct {
	Controller string  `json:"controller"`
	Sensor     string  `json:"sensor"`
	Value      float64 `json:"value"`
	Timestamp  string  `json:"timestamp"`
}

type SensorRange struct {
	SensorType *SensorType `json:"sensorType"`
	Min        float64     `json:"min"`
	Max        float64     `json:"max"`
}

type SensorType struct {
	Name string `json:"name"`
	Unit string `json:"unit"`
}

type sensorTypes struct {
	Types []*SensorType `json:"types"`
}

type sensorDataSet struct {
	SensorData []*SensorData `json:"data"`
}

type sensorDataPost struct {
	Data []*SensorData `json:"data"`
}
