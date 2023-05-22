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
