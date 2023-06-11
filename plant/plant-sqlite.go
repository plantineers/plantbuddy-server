// Author: Maximilian Floto, Yannick Kirschen
package plant

import (
	"context"
	"database/sql"
	"errors"

	"github.com/plantineers/plantbuddy-server/care_tips"
	"github.com/plantineers/plantbuddy-server/db"
)

// PlantSqliteRepository implements the PlantRepository interface.
// It uses a SQLite database as its data source.
type PlantSqliteRepository struct {
	db                   *sql.DB
	plantGroupRepository PlantGroupRepository
	careTipsRepository   care_tips.CareTipsRepository
}

// NewPlantRepository creates a new repository for plants.
// It will use the configured driver and data source from `buddy.json`
func NewPlantRepository(session *db.Session) (PlantRepository, error) {
	if !session.IsOpen() {
		return nil, errors.New("session is not open")
	}

	plantGroupRepository, err := NewPlantGroupRepository(session)
	if err != nil {
		return nil, err
	}

	careTipsRepository, err := care_tips.NewCareTipsRepository(session)
	if err != nil {
		return nil, err
	}

	return &PlantSqliteRepository{
		db:                   session.DB,
		plantGroupRepository: plantGroupRepository,
		careTipsRepository:   careTipsRepository,
	}, nil
}

func (r *PlantSqliteRepository) GetById(id int64) (*Plant, error) {
	var plantId int64
	var plantDescription *string
	var plantName *string
	var plantSpecies *string
	var plantLocation *string
	var plantGroupId int64

	err := r.db.QueryRow(`
    SELECT
        P.ID,
        P.PLANT_GROUP,
        P.DESCRIPTION,
        P.NAME,
        P.SPECIES,
        P.LOCATION
    FROM PLANT P
    WHERE P.ID = ?;`, id).Scan(&plantId, &plantGroupId, &plantDescription, &plantName, &plantSpecies, &plantLocation)

	if err != nil {
		return nil, err
	}

	if plantDescription == nil {
		plantDescription = new(string)
	}

	plantGroup, err := r.plantGroupRepository.GetById(plantGroupId)
	if err != nil {
		return nil, err
	}

	careTips, err := r.careTipsRepository.GetAdditionalByPlantId(plantId)
	if err != nil {
		return nil, err
	}

	if careTips == nil {
		careTips = make([]string, 0)
	}

	return &Plant{
		ID:                 plantId,
		Description:        *plantDescription,
		Name:               *plantName,
		Species:            *plantSpecies,
		Location:           *plantLocation,
		PlantGroup:         plantGroup,
		AdditionalCareTips: careTips,
	}, nil
}

func (r *PlantSqliteRepository) GetAll(filter *plantsFilter) ([]int64, error) {
	rows, err := r.getAllApplyFilter(filter)
	if err != nil {
		return nil, err
	}

	// Iterate over all rows and query the ID of the plant.
	var plantIds []int64
	for rows.Next() {
		var plantId int64

		err = rows.Scan(&plantId)
		if err != nil {
			return nil, err
		}

		plantIds = append(plantIds, plantId)
	}

	return plantIds, nil
}

func (r *PlantSqliteRepository) getAllApplyFilter(filter *plantsFilter) (*sql.Rows, error) {
	if filter != nil {
		return r.db.Query(`SELECT ID FROM PLANT WHERE PLANT_GROUP = ?;`, filter.PlantGroupId)
	}

	return r.db.Query(`SELECT ID FROM PLANT;`)
}

func (r *PlantSqliteRepository) Create(plant *plantChange) (*Plant, error) {
	tx, _ := r.db.BeginTx(context.Background(), nil)

	_, err := r.plantGroupRepository.GetById(plant.PlantGroupId)
	if err != nil {
		tx.Rollback()
		return nil, ErrPlantGroupNotExisting
	}

	plantStatement, err := r.db.Prepare(`
    INSERT INTO PLANT
        (PLANT_GROUP, DESCRIPTION, NAME, SPECIES, LOCATION)
    VALUES
        (?, ?, ?, ?, ?);`)

	result, err := plantStatement.Exec(
		plant.PlantGroupId,
		plant.Description,
		plant.Name,
		plant.Species,
		plant.Location,
	)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	plantId, _ := result.LastInsertId()

	err = r.careTipsRepository.CreateAdditionalByPlantId(plantId, plant.AdditionalCareTips)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return r.GetById(plantId)
}

func (r *PlantSqliteRepository) Update(id int64, plant *plantChange) (*Plant, error) {
	tx, _ := r.db.BeginTx(context.Background(), nil)

	_, err := r.plantGroupRepository.GetById(plant.PlantGroupId)
	if err != nil {
		tx.Rollback()
		return nil, ErrPlantGroupNotExisting
	}

	_, err = r.db.Exec(`
    UPDATE PLANT
    SET
        PLANT_GROUP = ?,
        DESCRIPTION = ?,
        NAME = ?,
        SPECIES = ?,
        LOCATION = ?
    WHERE ID = ?;`,
		plant.PlantGroupId,
		plant.Description,
		plant.Name,
		plant.Species,
		plant.Location,
		id)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = r.careTipsRepository.DeleteAdditionalByPlantId(id)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = r.careTipsRepository.CreateAdditionalByPlantId(id, plant.AdditionalCareTips)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return r.GetById(id)
}

func (r *PlantSqliteRepository) DeleteById(id int64) error {
	tx, _ := r.db.BeginTx(context.Background(), nil)

	_, err := r.db.Exec(`DELETE FROM PLANT WHERE ID = ?;`, id)

	if err != nil {
		tx.Rollback()
		return err
	}

	err = r.careTipsRepository.DeleteAdditionalByPlantId(id)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *PlantSqliteRepository) GetAllOverview() ([]PlantStub, error) {
	rows, err := r.db.Query(`
    SELECT
        P.ID,
        P.NAME
        FROM PLANT P;`)

	if err != nil {
		return nil, err
	}

	var plantStubs []PlantStub
	for rows.Next() {
		var plantId int64
		var plantName string

		err = rows.Scan(&plantId, &plantName)
		if err != nil {
			return nil, err
		}

		plantStubs = append(plantStubs, PlantStub{
			ID:   plantId,
			Name: plantName,
		})
	}

	return plantStubs, nil
}
