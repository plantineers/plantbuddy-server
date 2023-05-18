// Author: Yannick Kirschen
package sensor

type SensorDataFilter struct {
	Sensor string // Sensor Type
	Plant  int64  // Plant ID
	From   string // ISO 8601
	To     string // ISO 8601
}
