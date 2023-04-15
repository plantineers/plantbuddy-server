// Author: Maximilian Floto, Yannick Kirschen
package plant

import (
	"database/sql"
	"errors"
	"log"

	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/model"
)

type PlantSqliteRepository struct {
	db *sql.DB
}

// NewRepository creates a new repository for plants.
// It will use the configured driver and data source from `buddy.json`
func NewRepository(session *db.Session) (*PlantSqliteRepository, error) {
	if !session.IsOpen() {
		return nil, errors.New("session is not open")
	}

	return &PlantSqliteRepository{db: session.DB}, nil
}

// GetById returns a plant by its ID.
// If the plant does not exist, it will return nil.
func (r *PlantSqliteRepository) GetById(id int64) (*model.Plant, error) {
	var plantId int64
	var plantDescription *string
	var plantGroupId int64
	var plantGroupName string
	var plantGroupDescription *string
	var err = r.db.QueryRow(`
    SELECT P.ID AS PLANT_ID,
        P.DESCRIPTION AS PLANT_DESCRIPTION,
        PG.ID AS PLANT_GROUP_ID,
        PG.NAME AS PLANT_GROUP_NAME,
        PG.DESCRIPTION AS PLANT_GROUP_DESCRIPTION
    FROM PLANT P
    LEFT JOIN PLANT_GROUP PG on P.ID = PG.ID;`).Scan(
		&plantId,
		&plantDescription,
		&plantGroupId,
		&plantGroupName,
		&plantGroupDescription,
	)

	if err != nil {
		log.Fatal(err)
	}

	if plantDescription == nil {
		plantDescription = new(string)
	}

	if plantGroupDescription == nil {
		plantGroupDescription = new(string)
	}

	var plant = model.Plant{
		ID:          plantId,
		Description: *plantDescription,
		PlantGroup: &model.PlantGroup{
			ID:          plantGroupId,
			Name:        plantGroupName,
			Description: *plantGroupDescription,
		},
	}

	return &plant, nil
}
