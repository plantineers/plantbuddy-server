// Author: Yannick Kirschen
package sensor

import (
	"database/sql"
	"errors"
	"log"

	"github.com/plantineers/plantbuddy-server/db"
)

// SensorTypeRepository provides access to sensor types.
// It uses a SQLite database as data source.
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

func (r *SensorTypeSqliteRepository) GetAll() ([]*SensorType, error) {
	var rows, err = r.db.Query(`SELECT NAME, UNIT FROM SENSOR_TYPE;`)

	if err != nil {
		log.Fatal(err)
	}

	var types []*SensorType
	for rows.Next() {
		var name string
		var unit string

		err = rows.Scan(&name, &unit)
		if err != nil {
			rows.Close()
			return nil, err
		}
		types = append(types, &SensorType{Name: name, Unit: unit})
	}

	return types, nil
}
