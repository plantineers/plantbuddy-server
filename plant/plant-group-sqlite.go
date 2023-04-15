package plant

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/model"
)

type PlantGroupSqliteRepository struct {
	db *sql.DB
}

// NewRepository creates a new repository for plant-groups.
// It will use the configured driver and data source from `buddy.json`
func NewPlantGroupRepository(session *db.Session) (PlantGroupRepository, error) {
	if !session.IsOpen() {
		return nil, errors.New("session is not open")
	}

	return &PlantGroupSqliteRepository{db: session.DB}, nil
}

// GetPlantGroupById returns a plant group by its ID.
func (r *PlantGroupSqliteRepository) GetById(id int64) (*model.PlantGroup, error) {
	var plantGroupId int64
	var plantGroupName string
	var plantGroupDescription *string
	err := r.db.QueryRow(`
    SELECT PG.ID AS PLANT_GROUP_ID,
        PG.NAME AS PLANT_GROUP_NAME,
        PG.DESCRIPTION AS PLANT_GROUP_DESCRIPTION FROM PLANT_GROUP PG
    WHERE PG.ID = ?;`, id).Scan(&plantGroupId, &plantGroupName, &plantGroupDescription)

	if err != nil {
		fmt.Print(err)
		return nil, err
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
func (r *PlantGroupSqliteRepository) GetAll() (*model.PlantGroups, error) {
	var plantGroupIds []int64
	rows, err := r.db.Query(`SELECT ID FROM PLANT_GROUP;`)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	// Iterate over all rows and query the plant group by its ID.
	for rows.Next() {
		var plantGroupId int64
		err = rows.Scan(&plantGroupId)
		if err != nil {
			fmt.Print(err)
			return nil, err
		}
		plantGroupIds = append(plantGroupIds, plantGroupId)
	}
	plantGroups := model.PlantGroups{
		PlantGroups: plantGroupIds,
	}
	return &plantGroups, nil
}
