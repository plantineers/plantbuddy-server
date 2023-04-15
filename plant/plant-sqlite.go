package plant

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/plantineers/plantbuddy-server/model"
)

type PlantSqliteRepository struct {
	db *sql.DB // TODO: make this a connection pool
}

func (r *PlantSqliteRepository) GetById(id int64) (*model.Plant, error) {
	db, err := sql.Open("sqlite3", "./plantbuddy.sqlite")

	if err != nil {
		return nil, err
	}

	r.db = db

	defer db.Close()

	var plantId int64
	var plantDescription *string
	var plantGroupId int64
	var plantGroupName string
	var plantGroupDescription *string
	err = db.QueryRow(`
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
