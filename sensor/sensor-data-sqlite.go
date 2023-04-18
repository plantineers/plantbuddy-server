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
	db               *sql.DB
	sensorRepository SensorRepository
}

// SensorDataSqliteRepository creates a new repository for sensor data.
// It will use the configured driver and data source from `buddy.json`
func NewSensorDataRepository(session *db.Session) (SensorDataRepository, error) {
	if !session.IsOpen() {
		return nil, errors.New("session is not open")
	}

	sensorRepository, err := NewSensorRepository(session)
	if err != nil {
		return nil, err
	}

	return &SensorDataSqliteRepository{
		db:               session.DB,
		sensorRepository: sensorRepository,
	}, nil
}

func (r *SensorDataSqliteRepository) GetAll(filter *SensorDataFilter) ([]*model.SensorData, error) {
	rows, err := r.db.Query(`
    SELECT SD.ID, S.ID, SD.VALUE, SD.TIMESTAMP
    FROM SENSOR_DATA SD
    JOIN SENSOR S ON SD.SENSOR = S.ID
    WHERE S.ID = 1
        AND S.PLANT = 1
        AND SD.TIMESTAMP BETWEEN DATETIME(?) AND DATETIME(?);`, filter.From, filter.To)
	if err != nil {
		return nil, err
	}

	var data []*model.SensorData
	for rows.Next() {
		var id int64
		var sensor int64
		var value float64
		var timestamp string

		err = rows.Scan(&id, &sensor, &value, &timestamp)
		if err != nil {
			return nil, err
		}

		data = append(data, &model.SensorData{
			ID:        id,
			Sensor:    sensor,
			Value:     value,
			Timestamp: timestamp,
		})
	}

	return data, nil
}

func (r *SensorDataSqliteRepository) Save(data *model.SensorData) error {
	tx, _ := r.db.BeginTx(context.Background(), nil)

	sensor, err := r.sensorRepository.GetById(data.Sensor)
	if err == sql.ErrNoRows {
		return errors.New("sensor does not exist")
	} else if err != nil {
		return err
	}

	statement, err := r.db.Prepare("INSERT INTO SENSOR_DATA (SENSOR, VALUE, TIMESTAMP) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer statement.Close()

	res, err := statement.Exec(sensor.ID, data.Value, data.Timestamp)
	if err != nil {
		tx.Rollback()
		return err
	}

	id, _ := res.LastInsertId()
	err = tx.Commit()
	if err != nil {
		return err
	}
	data.ID = id

	return nil
}
