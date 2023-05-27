// Author: Yannick Kirschen
package sensor_range

import (
	"database/sql"
	"errors"

	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/sensor"
)

type SensorRangeSqliteRepository struct {
	db *sql.DB
}

// NewCareTipsRepository creates a new repository for care tips.
// It will use the configured driver and data source from `buddy.json`
func NewSensorRangeRepository(session *db.Session) (SensorRangeRepository, error) {
	if !session.IsOpen() {
		return nil, errors.New("session is not open")
	}

	return &SensorRangeSqliteRepository{db: session.DB}, nil
}

func (r *SensorRangeSqliteRepository) GetAllByPlantGroupId(id int64) ([]*sensor.SensorRange, error) {
	var sensorRanges []*sensor.SensorRange
	rows, err := r.db.Query(`
    SELECT SR.MIN, SR.MAX, ST.NAME, ST.UNIT
        FROM SENSOR_RANGE SR
        LEFT JOIN SENSOR_TYPE ST on SR.SENSOR = ST.NAME
        WHERE SR.PLANT_GROUP = ?;`, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var sensorRange sensor.SensorRange
		var sensorType sensor.SensorType

		err := rows.Scan(&sensorRange.Min, &sensorRange.Max, &sensorType.Name, &sensorType.Unit)
		if err != nil {
			return nil, err
		}
		sensorRange.SensorType = &sensorType
		sensorRanges = append(sensorRanges, &sensorRange)
	}

	return sensorRanges, nil
}

func (r *SensorRangeSqliteRepository) Create(plantGroupId int64, sensorRange *sensor.SensorRangeChange) error {
	_, err := r.db.Exec(`
    INSERT INTO SENSOR_RANGE (PLANT_GROUP, SENSOR, MIN, MAX)
        VALUES (?, ?, ?, ?);`, plantGroupId, sensorRange.Sensor, sensorRange.Min, sensorRange.Max)

	return err
}

func (r *SensorRangeSqliteRepository) CreateAll(plantGroupId int64, sensorRanges []*sensor.SensorRangeChange) error {
	for _, sensorRange := range sensorRanges {
		err := r.Create(plantGroupId, sensorRange)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *SensorRangeSqliteRepository) DeleteAllByPlantGroupId(id int64) error {
	_, err := r.db.Exec(`DELETE FROM SENSOR_RANGE WHERE PLANT_GROUP = ?;`, id)
	return err
}
