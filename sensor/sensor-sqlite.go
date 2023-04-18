package sensor

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/model"
)

type SensorSqliteRepository struct {
	db *sql.DB
}

// NewSensorRepository creates a new repository for sensors.
// It will use the configured driver and data source from `buddy.json`
func NewSensorRepository(session *db.Session) (SensorRepository, error) {
	if !session.IsOpen() {
		return nil, errors.New("session is not open")
	}

	return &SensorSqliteRepository{db: session.DB}, nil
}

func (r *SensorSqliteRepository) GetById(id int64) (*model.Sensor, error) {
	var sensorId int64
	var plantId int64
	var interval int64
	var sensorTypeId int64
	var sensorTypeName string
	var sensorTypeUnit string
	var err = r.db.QueryRow(`
    SELECT S.ID AS SENSOR_ID,
       S.PLANT AS PLANT_ID,
       S.INTERVAL AS INTERVAL,
       ST.ID AS SENSOR_TYPE_ID,
       ST.NAME AS SENSOR_TYPE_NAME,
       ST.UNIT AS SENSOR_TYPE_UNIT
    FROM SENSOR S
    LEFT JOIN SENSOR_TYPE ST
    WHERE SENSOR_ID = ?;`, id).Scan(
		&sensorId,
		&plantId,
		&interval,
		&sensorTypeId,
		&sensorTypeName,
		&sensorTypeUnit,
	)

	if err != nil {
		return nil, err
	}

	var sensor = model.Sensor{
		ID:       sensorId,
		Plant:    plantId,
		Interval: interval,
		SensorType: &model.SensorType{
			Name: sensorTypeName,
			Unit: sensorTypeUnit,
		},
	}

	return &sensor, nil
}

func (r *SensorSqliteRepository) GetAllIds() ([]int64, error) {
	var rows, err = r.db.Query(`SELECT ID FROM SENSOR;`)

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

func (r *SensorSqliteRepository) Create(sensor *sensorChange) (*model.Sensor, error) {

	tx, _ := r.db.BeginTx(context.Background(), nil)

	if r.db.QueryRow(`SELECT ID FROM SENSOR_TYPE WHERE ID = ?;`, sensor.SensorType).Scan(&sensor.SensorType) != nil {
		return nil, errors.New("Sensor type does not exist")
	}

	if r.db.QueryRow(`SELECT ID FROM PLANT WHERE ID = ?;`, sensor.Plant).Scan(&sensor.Plant) != nil {
		return nil, errors.New("Plant does not exist")
	}

	var statement, err = r.db.Prepare(`
    INSERT INTO SENSOR (PLANT, TYPE, INTERVAL)
    VALUES (?, ?, ?);`)
	defer statement.Close()

	if err != nil {
		return nil, err
	}

	result, err := statement.Exec(
		sensor.Plant,
		sensor.SensorType,
		sensor.Interval,
	)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	id, _ := result.LastInsertId()
	createdSensor, _ := r.GetById(id)

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return createdSensor, err
}

func (r *SensorSqliteRepository) Update(sensor *sensorChange, id int64) (*model.Sensor, error) {
	var statement, err = r.db.Prepare(`
    UPDATE SENSOR
    SET PLANT = ?,
        TYPE = ?,
        INTERVAL = ?
    WHERE ID = ?;`)
	defer statement.Close()

	if err != nil {
		return nil, err
	}

	_, err = statement.Exec(
		sensor.Plant,
		sensor.SensorType,
		sensor.Interval,
		id,
	)

	if err != nil {
		return nil, err
	}

	updatedSensor, _ := r.GetById(id)

	return updatedSensor, err
}

func (r *SensorSqliteRepository) Delete(id int64) error {
	var statement, err = r.db.Prepare(`DELETE FROM SENSOR WHERE ID = ?;`)
	defer statement.Close()

	if err != nil {
		return err
	}

	_, err = statement.Exec(id)

	return err
}
