// Author: Yannick Kirschen
package controller

type ControllerRepository interface {
	GetAllUUIDs() ([]string, error)

	GetByUUID(uuid string) (*Controller, error)
}
