package plant

import "github.com/plantineers/plantbuddy-server/model"

type PlantRepository interface {
	GetById(id int64) (*model.Plant, error)
	GetAll() (*[]*model.Plant, error)
}

type PlantGroupRepository interface {
	GetById(id int64) (*model.PlantGroup, error)
	GetAll() (*[]*model.PlantGroup, error)
}
