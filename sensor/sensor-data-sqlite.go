// Author: Yannick Kirschen
package sensor

import (
	"context"
	"database/sql"
	"errors"

	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/model"
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

func (r *SensorDataSqliteRepository) GetAll(filter *SensorDataFilter) ([]*model.SensorData, error) {
	rows, err := r.db.Query(`
    SELECT SD.CONTROLLER,
       SD.SENSOR,
       SD.VALUE,
       SD.TIMESTAMP
    FROM SENSOR_DATA SD
    LEFT JOIN CONTROLLER C on SD.CONTROLLER = C.UUID
    WHERE C.PLANT = ?
        AND SD.SENSOR = ?
        AND SD.TIMESTAMP BETWEEN DATETIME(?) AND DATETIME(?);`,
		filter.Plant, filter.Sensor, filter.From, filter.To)
	if err != nil {
		return nil, err
	}

	var data []*model.SensorData
	for rows.Next() {
		var controller string
		var sensor string
		var value float64
		var timestamp string

		err = rows.Scan(&controller, &sensor, &value, &timestamp)
		if err != nil {
			return nil, err
		}

		data = append(data, &model.SensorData{
			Controller: controller,
			Sensor:     sensor,
			Value:      value,
			Timestamp:  timestamp,
		})
	}

	return data, nil
}

func (r *SensorDataSqliteRepository) Save(data *model.SensorData) error {
	tx, _ := r.db.BeginTx(context.Background(), nil)

	statement, err := r.db.Prepare("INSERT INTO SENSOR_DATA (CONTROLLER, SENSOR, VALUE, TIMESTAMP) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(data.Controller, data.Sensor, data.Value, data.Timestamp)
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

func (r *SensorDataSqliteRepository) SaveAll(data []*model.SensorData) []error {
	var errors []error
	for _, d := range data {
		err := r.Save(d)
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}
