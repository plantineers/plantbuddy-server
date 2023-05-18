package sensor_range

import (
	"database/sql"
	"errors"

	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/model"
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

func (r *SensorRangeSqliteRepository) GetAllByPlantGroupId(id int64) ([]*model.SensorRange, error) {
	var sensorRanges []*model.SensorRange
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
		var sensorRange model.SensorRange
		var sensorType model.SensorType

		err := rows.Scan(&sensorRange.Min, &sensorRange.Max, &sensorType.Name, &sensorType.Unit)
		if err != nil {
			return nil, err
		}
		sensorRange.SensorType = &sensorType
		sensorRanges = append(sensorRanges, &sensorRange)
	}

	return sensorRanges, nil
}
