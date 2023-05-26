// Author: Yannick Kirschen
package controller

import "github.com/plantineers/plantbuddy-server/model"

type ControllerRepository interface {
	GetAllUUIDs() ([]string, error)

	GetByUUID(uuid string) (*model.Controller, error)
}
