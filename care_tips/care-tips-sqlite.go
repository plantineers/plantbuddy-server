// Author: Maximilian Floto, Yannick Kirschen
package care_tips

import (
	"database/sql"
	"errors"

	"github.com/plantineers/plantbuddy-server/db"
)

// CareTipsSqliteRepository implements the CareTipsRepository interface.
// It uses a SQLite database as its data source.
type CareTipsSqliteRepository struct {
	db *sql.DB
}

// NewCareTipsRepository creates a new repository for care tips.
// It will use the configured driver and data source from `buddy.json`
func NewCareTipsRepository(session *db.Session) (CareTipsRepository, error) {
	if !session.IsOpen() {
		return nil, errors.New("session is not open")
	}

	return &CareTipsSqliteRepository{db: session.DB}, nil
}

func (r *CareTipsSqliteRepository) GetByPlantGroupId(id int64) ([]string, error) {
	rows, err := r.db.Query(`SELECT CT.TIP FROM CARE_TIPS CT WHERE CT.PLANT_GROUP = ?;`, id)

	if err != nil {
		return nil, err
	}

	var careTips []string
	for rows.Next() {
		var careTip string

		err = rows.Scan(&careTip)
		if err != nil {
			return nil, err
		}

		careTips = append(careTips, careTip)
	}

	return careTips, nil
}

func (r *CareTipsSqliteRepository) GetAdditionalByPlantId(id int64) ([]string, error) {
	rows, err := r.db.Query(`SELECT CT.TIP FROM ADDITIONAL_CARE_TIPS CT WHERE CT.PLANT = ?;`, id)

	if err != nil {
		return nil, err
	}

	var careTips []string
	for rows.Next() {
		var careTip string

		err = rows.Scan(&careTip)
		if err != nil {
			return nil, err
		}

		careTips = append(careTips, careTip)
	}

	return careTips, nil
}

func (r *CareTipsSqliteRepository) Create(plantGroupId int64, careTips []string) error {
	for _, careTip := range careTips {
		_, err := r.db.Exec(`
        INSERT INTO CARE_TIPS (PLANT_GROUP, TIP)
        VALUES (?, ?);`, plantGroupId, careTip)

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *CareTipsSqliteRepository) DeleteAllByPlantGroupId(id int64) error {
	_, err := r.db.Exec(`DELETE FROM CARE_TIPS WHERE PLANT_GROUP = ?;`, id)
	return err
}

func (r *CareTipsSqliteRepository) CreateAdditionalByPlantId(plantId int64, careTips []string) error {
	for _, careTip := range careTips {
		_, err := r.db.Exec(`INSERT INTO ADDITIONAL_CARE_TIPS (PLANT, TIP) VALUES (?, ?);`, plantId, careTip)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *CareTipsSqliteRepository) DeleteAdditionalByPlantId(plantId int64) error {
	_, err := r.db.Exec(`DELETE FROM ADDITIONAL_CARE_TIPS WHERE PLANT = ?;`, plantId)
	if err != nil {
		return err
	}

	return nil
}
