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

// GetPlantById returns a plant by its ID.
// If the plant does not exist, it will return nil.
func (r *PlantSqliteRepository) GetPlantById(id int64) (*model.Plant, error) {
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
        PG.DESCRIPTION AS PLANT_GROUP_DESCRIPTION FROM PLANT P
    LEFT JOIN PLANT_GROUP PG on P.ID = PG.ID;`).Scan(&plantId, &plantDescription, &plantGroupId, &plantGroupName, &plantGroupDescription)

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

// Reads all plantIds from the database and returns them as a slice of plants.
func (r *PlantSqliteRepository) GetAllPlants() (*[]*model.Plant, error) {
	var plants []*model.Plant
	rows, err := r.db.Query(`SELECT ID FROM PLANT;`)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Iterate over all rows and query the plant by its ID.
	for rows.Next() {
		var plantId int64
		err = rows.Scan(&plantId)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		plant, err := r.GetPlantById(plantId)
		if err != nil {
			return nil, err
		}
		plants = append(plants, plant)
	}
	return &plants, nil
}

// GetPlantGroupById returns a plant group by its ID.
func (r *PlantSqliteRepository) GetPlantGroupById(id int64) (*model.PlantGroup, error) {
	var plantGroupId int64
	var plantGroupName string
	var plantGroupDescription *string
	err := r.db.QueryRow(`
    SELECT PG.ID AS PLANT_GROUP_ID,
        PG.NAME AS PLANT_GROUP_NAME,
        PG.DESCRIPTION AS PLANT_GROUP_DESCRIPTION FROM PLANT_GROUP PG
    WHERE PG.ID = ?;`, id).Scan(&plantGroupId, &plantGroupName, &plantGroupDescription)

	if err != nil {
		log.Fatal(err)
	}

	if plantGroupDescription == nil {
		plantGroupDescription = new(string)
	}

	var plantGroup = model.PlantGroup{
		ID:          plantGroupId,
		Name:        plantGroupName,
		Description: *plantGroupDescription,
	}

	return &plantGroup, nil
}

// Reads all plantGroupIds from the database and returns them as a slice of plant groups.
func (r *PlantSqliteRepository) GetAllPlantGroups() (*[]*model.PlantGroup, error) {
	var plantGroups []*model.PlantGroup
	rows, err := r.db.Query(`SELECT ID FROM PLANT_GROUP;`)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Iterate over all rows and query the plant group by its ID.
	for rows.Next() {
		var plantGroupId int64
		err = rows.Scan(&plantGroupId)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		plantGroup, err := r.GetPlantGroupById(plantGroupId)
		if err != nil {
			return nil, err
		}
		plantGroups = append(plantGroups, plantGroup)
	}
	return &plantGroups, nil
}
