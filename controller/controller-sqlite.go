// Author: Yannick Kirschen
package controller

import (
	"database/sql"
	"errors"

	"github.com/plantineers/plantbuddy-server/db"
)

type ControllerSqliteRepository struct {
	db *sql.DB
}

// NewControllerRepository creates a new repository for care tips.
// It will use the configured driver and data source from `buddy.json`
func NewControllerRepository(session *db.Session) (ControllerRepository, error) {
	if !session.IsOpen() {
		return nil, errors.New("session is not open")
	}

	return &ControllerSqliteRepository{db: session.DB}, nil
}

func (r *ControllerSqliteRepository) GetAllUUIDs() ([]string, error) {
	rows, err := r.db.Query(`SELECT C.UUID FROM CONTROLLER C;`)

	if err != nil {
		return nil, err
	}

	var uuids []string
	for rows.Next() {
		var uuid string

		err = rows.Scan(&uuid)
		if err != nil {
			return nil, err
		}

		uuids = append(uuids, uuid)
	}

	return uuids, nil
}

func (r *ControllerSqliteRepository) GetByUUID(uuid string) (*Controller, error) {
	var controller Controller

	err := r.db.QueryRow(`
    SELECT C.UUID, C.PLANT_GROUP
        FROM CONTROLLER C
        WHERE C.UUID = ?;`, uuid).Scan(&controller.UUID, &controller.PlantGroup)

	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(`
    SELECT DISTINCT SD.SENSOR
        FROM SENSOR_DATA SD
        WHERE SD.CONTROLLER = ?;`, uuid)

	if err != nil {
		return nil, err
	}

	var sensors []string
	for rows.Next() {
		var sensor string

		err = rows.Scan(&sensor)
		if err != nil {
			return nil, err
		}

		sensors = append(sensors, sensor)
	}

	controller.Sensors = sensors
	return &controller, nil
}
