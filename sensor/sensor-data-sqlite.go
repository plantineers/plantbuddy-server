// Author: Yannick Kirschen
package sensor

import (
	"context"
	"database/sql"
	"errors"

	"github.com/plantineers/plantbuddy-server/db"
)

type SensorDataSqliteRepository struct {
	db *sql.DB
}

// SensorDataSqliteRepository creates a new repository for sensor data.
// It will use the configured driver and data source from `buddy.json`
func NewSensorDataRepository(session *db.Session) (SensorDataRepository, error) {
	if !session.IsOpen() {
		return nil, errors.New("session is not open")
	}

	return &SensorDataSqliteRepository{db: session.DB}, nil
}

func (r *SensorDataSqliteRepository) GetAll(filter *SensorDataFilter) ([]*SensorData, error) {
	var plantGroupId int64

	if filter.Plant != 0 {
		err := r.db.QueryRow(`
        SELECT
            P.PLANT_GROUP
        FROM PLANT P
        WHERE P.ID = ?;`, filter.Plant).Scan(&plantGroupId)
		if err != nil {
			return nil, err
		}
	} else {
		plantGroupId = filter.PlantGroup
	}

	rows, err := r.db.Query(`
    SELECT SD.CONTROLLER,
       SD.SENSOR,
       SD.VALUE,
       SD.TIMESTAMP
    FROM SENSOR_DATA SD
    LEFT JOIN CONTROLLER C on SD.CONTROLLER = C.UUID
    WHERE C.PLANT_GROUP = ?
        AND SD.SENSOR = ?
        AND SD.TIMESTAMP BETWEEN DATETIME(?) AND DATETIME(?);`,
		plantGroupId, filter.Sensor, filter.From, filter.To)
	if err != nil {
		return nil, err
	}

	var data []*SensorData
	for rows.Next() {
		var controller string
		var sensor string
		var value float64
		var timestamp string

		err = rows.Scan(&controller, &sensor, &value, &timestamp)
		if err != nil {
			return nil, err
		}

		data = append(data, &SensorData{
			Controller: controller,
			Sensor:     sensor,
			Value:      value,
			Timestamp:  timestamp,
		})
	}

	return data, nil
}

func (r *SensorDataSqliteRepository) Save(data *SensorData) error {
	tx, _ := r.db.BeginTx(context.Background(), nil)

	_, err := r.db.Exec("INSERT INTO SENSOR_DATA (CONTROLLER, SENSOR, VALUE, TIMESTAMP) VALUES (?, ?, ?, ?)",
		data.Controller, data.Sensor, data.Value, data.Timestamp)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *SensorDataSqliteRepository) SaveAll(data []*SensorData) []error {
	var errors []error
	for _, d := range data {
		err := r.Save(d)
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}
