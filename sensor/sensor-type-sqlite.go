package sensor

import (
	"database/sql"
	"errors"
	"log"

	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/model"
)

type SensorTypeSqliteRepository struct {
	db *sql.DB
}

// NewSensorTypeRepository creates a new repository for sensor types.
// It will use the configured driver and data source from `buddy.json`
func NewSensorTypeRepository(session *db.Session) (SensorTypeRepository, error) {
	if !session.IsOpen() {
		return nil, errors.New("session is not open")
	}

	return &SensorTypeSqliteRepository{db: session.DB}, nil
}

func (r *SensorTypeSqliteRepository) GetById(id int64) (*model.SensorType, error) {
	var name string
	var unit string
	var err = r.db.QueryRow(`
    SELECT NAME, UNIT
    FROM SENSOR_TYPE
    WHERE ID = ?;`, id).Scan(&name, &unit)

	if err != nil {
		return nil, err
	}

	return &model.SensorType{Name: name, Unit: unit}, nil
}

func (r *SensorTypeSqliteRepository) GetAllIds() ([]int64, error) {
	var rows, err = r.db.Query(`SELECT ID FROM SENSOR_TYPE;`)

	if err != nil {
		log.Fatal(err)
	}

	var ids []int64
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			rows.Close()
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}
