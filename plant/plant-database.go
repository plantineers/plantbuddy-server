// Author: Maximilian Floto, Yannick Kirschen
package plant

import "github.com/plantineers/plantbuddy-server/model"

type PlantRepository interface {
	GetById(id int64) (*model.Plant, error)
}
