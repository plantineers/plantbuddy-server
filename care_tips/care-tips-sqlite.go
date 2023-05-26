package care_tips

import (
	"database/sql"
	"errors"

	"github.com/plantineers/plantbuddy-server/db"
)

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
	statement, err := r.db.Prepare(`
        INSERT INTO CARE_TIPS (PLANT_GROUP, TIP)
        VALUES (?, ?);`)

	if err != nil {
		return err
	}

	for _, careTip := range careTips {
		_, err = statement.Exec(plantGroupId, careTip)
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
