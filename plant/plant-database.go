package plant

import "github.com/plantineers/plantbuddy-server/model"

type PlantRepository interface {
	GetPlantById(id int64) (*model.Plant, error)
	GetAllPlants() ([]*model.Plant, error)
}

type PlantGroupRepository interface {
	GetPlantGroupById(id int64) (*model.PlantGroup, error)
	GetAllPlantGroups() ([]*model.PlantGroup, error)
}
