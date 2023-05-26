package plant

import (
	"context"
	"database/sql"
	"errors"

	"github.com/plantineers/plantbuddy-server/care_tips"
	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/model"
	"github.com/plantineers/plantbuddy-server/sensor_range"
)

type PlantGroupSqliteRepository struct {
	db                    *sql.DB
	careTipsRepository    care_tips.CareTipsRepository
	sensorRangeRepository sensor_range.SensorRangeRepository
}

// NewPlantGroupRepository creates a new repository for plant-groups.
// It will use the configured driver and data source from `buddy.json`
func NewPlantGroupRepository(session *db.Session) (PlantGroupRepository, error) {
	if !session.IsOpen() {
		return nil, errors.New("session is not open")
	}

	careTipsRepository, err := care_tips.NewCareTipsRepository(session)
	if err != nil {
		return nil, err
	}

	sensorRangeRepository, err := sensor_range.NewSensorRangeRepository(session)
	if err != nil {
		return nil, err
	}

	return &PlantGroupSqliteRepository{
		db:                    session.DB,
		careTipsRepository:    careTipsRepository,
		sensorRangeRepository: sensorRangeRepository,
	}, nil
}

func (r *PlantGroupSqliteRepository) GetById(id int64) (*model.PlantGroup, error) {
	var plantGroupId int64
	var plantGroupName string
	var plantGroupDescription *string
	err := r.db.QueryRow(`
    SELECT PG.ID,
        PG.NAME,
        PG.DESCRIPTION
        FROM PLANT_GROUP PG
    WHERE PG.ID = ?;`, id).Scan(&plantGroupId, &plantGroupName, &plantGroupDescription)

	if err != nil {
		return nil, err
	}

	if plantGroupDescription == nil {
		plantGroupDescription = new(string)
	}

	careTips, err := r.careTipsRepository.GetByPlantGroupId(plantGroupId)
	if err != nil {
		return nil, err
	}

	sensor_ranges, err := r.sensorRangeRepository.GetAllByPlantGroupId(plantGroupId)
	if err != nil {
		return nil, err
	}

	var plantGroup = model.PlantGroup{
		ID:           plantGroupId,
		Name:         plantGroupName,
		Description:  *plantGroupDescription,
		CareTips:     careTips,
		SensorRanges: sensor_ranges,
	}

	return &plantGroup, nil
}

func (r *PlantGroupSqliteRepository) GetAll() ([]int64, error) {
	var plantGroupIds []int64
	rows, err := r.db.Query(`SELECT ID FROM PLANT_GROUP;`)
	if err != nil {
		return nil, err
	}

	// Iterate over all rows and query the ID of the plant group..
	for rows.Next() {
		var plantGroupId int64
		err = rows.Scan(&plantGroupId)
		if err != nil {
			return nil, err
		}
		plantGroupIds = append(plantGroupIds, plantGroupId)
	}

	return plantGroupIds, nil
}

func (r *PlantGroupSqliteRepository) Create(plantGroup *plantGroupChange) (*model.PlantGroup, error) {
	tx, _ := r.db.BeginTx(context.Background(), nil)

	statement, err := r.db.Prepare(`
        INSERT INTO PLANT_GROUP (NAME, DESCRIPTION)
        VALUES (?, ?);`)

	if err != nil {
		return nil, err
	}

	result, err := statement.Exec(plantGroup.Name, plantGroup.Description)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	id, _ := result.LastInsertId()
	err = r.careTipsRepository.Create(id, plantGroup.CareTips)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = r.sensorRangeRepository.CreateAll(id, plantGroup.SensorRanges)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return r.GetById(id)
}

func (r *PlantGroupSqliteRepository) Update(id int64, plantGroup *plantGroupChange) (*model.PlantGroup, error) {
	tx, _ := r.db.BeginTx(context.Background(), nil)

	var statement, err = r.db.Prepare(`
        UPDATE PLANT_GROUP
        SET NAME = ?,
            DESCRIPTION = ?
        WHERE ID = ?;`,
	)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = statement.Exec(
		plantGroup.Name,
		plantGroup.Description,
		id,
	)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = r.careTipsRepository.DeleteAllByPlantGroupId(id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = r.careTipsRepository.Create(id, plantGroup.CareTips)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = r.sensorRangeRepository.UpdateAll(id, plantGroup.SensorRanges)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return r.GetById(id)
}

func (r *PlantGroupSqliteRepository) Delete(id int64) error {
	tx, _ := r.db.BeginTx(context.Background(), nil)

	var statement, err = r.db.Prepare(`DELETE FROM PLANT_GROUP WHERE ID = ?;`)

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = statement.Exec(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = r.careTipsRepository.DeleteAllByPlantGroupId(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = r.sensorRangeRepository.DeleteAllByPlantGroupId(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
