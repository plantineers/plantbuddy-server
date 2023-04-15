package plant

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/model"
)

type PlantSqliteRepository struct {
	db *sql.DB
}

// NewPlantRepository creates a new repository for plants.
// It will use the configured driver and data source from `buddy.json`
func NewPlantRepository(session *db.Session) (PlantRepository, error) {
	if !session.IsOpen() {
		return nil, errors.New("session is not open")
	}

	return &PlantSqliteRepository{db: session.DB}, nil
}

// GetPlantById returns a plant by its ID.
// If the plant does not exist, it will return nil.
func (r *PlantSqliteRepository) GetById(id int64) (*model.Plant, error) {
	var plantId int64
	var plantDescription *string
	var plantGroupId int64
	var plantGroupName string
	var plantGroupDescription *string
	err := r.db.QueryRow(`
    SELECT P.ID AS PLANT_ID,
        P.DESCRIPTION AS PLANT_DESCRIPTION,
        PG.ID AS PLANT_GROUP_ID,
        PG.NAME AS PLANT_GROUP_NAME,
        PG.DESCRIPTION AS PLANT_GROUP_DESCRIPTION FROM PLANT P
    LEFT JOIN PLANT_GROUP PG ON P.ID = PG.ID
        WHERE PLANT_ID = ?;`, id).Scan(&plantId, &plantDescription, &plantGroupId, &plantGroupName, &plantGroupDescription)

	if err != nil {
		fmt.Print(err)
		return nil, err
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

// Reads all plantIds from the database and returns them as a slice of plants.
func (r *PlantSqliteRepository) GetAll() ([]int64, error) {
	var plantIds []int64
	rows, err := r.db.Query(`SELECT ID FROM PLANT;`)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Iterate over all rows and query the ID of the plant.
	for rows.Next() {
		var plantId int64
		err = rows.Scan(&plantId)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		plantIds = append(plantIds, plantId)
	}

	return plantIds, nil
}
