package care_tips

type CareTipsRepository interface {
	GetByPlantGroupId(id int64) ([]string, error)

	GetAdditionalByPlantId(id int64) ([]string, error)

	CreateAdditionalByPlantId(plantId int64, careTips []string) error
}
