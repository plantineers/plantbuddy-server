// Author: Maximilian Floto, Yannick Kirschen
package plant

import (
	"database/sql"
	"errors"

	"github.com/plantineers/plantbuddy-server/care_tips"
	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/model"
)

type PlantSqliteRepository struct {
	db                   *sql.DB
	plantGroupRepository PlantGroupRepository
	careTipsRepository   care_tips.CareTipsRepository
}

// NewPlantRepository creates a new repository for plants.
// It will use the configured driver and data source from `buddy.json`
func NewPlantRepository(session *db.Session) (PlantRepository, error) {
	if !session.IsOpen() {
		return nil, errors.New("session is not open")
	}

	plantGroupRepository, err := NewPlantGroupRepository(session)
	if err != nil {
		return nil, err
	}

	careTipsRepository, err := care_tips.NewCareTipsRepository(session)
	if err != nil {
		return nil, err
	}

	return &PlantSqliteRepository{
		db:                   session.DB,
		plantGroupRepository: plantGroupRepository,
		careTipsRepository:   careTipsRepository,
	}, nil
}

func (r *PlantSqliteRepository) GetById(id int64) (*model.Plant, error) {
	var plantId int64
	var plantDescription *string
	var plantGroupId int64

	err := r.db.QueryRow(`
    SELECT
        P.ID,
        P.PLANT_GROUP,
        P.DESCRIPTION
    FROM PLANT P
    WHERE P.ID = ?;`, id).Scan(&plantId, &plantGroupId, &plantDescription)

	if err != nil {
		return nil, err
	}

	if plantDescription == nil {
		plantDescription = new(string)
	}

	plantGroup, err := r.plantGroupRepository.GetById(plantGroupId)
	if err != nil {
		return nil, err
	}

	careTips, err := r.careTipsRepository.GetAdditionalByPlantId(plantId)
	if err != nil {
		return nil, err
	}

	return &model.Plant{
		ID:                 plantId,
		Description:        *plantDescription,
		PlantGroup:         plantGroup,
		AdditionalCareTips: careTips,
	}, nil
}

func (r *PlantSqliteRepository) GetAll(filter *PlantsFilter) ([]int64, error) {
	rows, err := r.getAllApplyFilter(filter)
	if err != nil {
		return nil, err
	}

	// Iterate over all rows and query the ID of the plant.
	var plantIds []int64
	for rows.Next() {
		var plantId int64

		err = rows.Scan(&plantId)
		if err != nil {
			return nil, err
		}

		plantIds = append(plantIds, plantId)
	}

	return plantIds, nil
}

func (r *PlantSqliteRepository) getAllApplyFilter(filter *PlantsFilter) (*sql.Rows, error) {
	if filter != nil {
		return r.db.Query(`SELECT ID FROM PLANT WHERE PLANT_GROUP = ?;`, filter.PlantGroupId)
	}

	return r.db.Query(`SELECT ID FROM PLANT;`)
}
