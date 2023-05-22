package user_management

import (
	"database/sql"
	"errors"
	"github.com/plantineers/plantbuddy-server/db"
	"github.com/plantineers/plantbuddy-server/model"
)

type UserSqliteRepository struct {
	db *sql.DB
}

func NewUserRepository(session *db.Session) (UserRepository, error) {
	if !session.IsOpen() {
		return nil, errors.New("session is not open")
	}

	return &UserSqliteRepository{
		db: session.DB,
	}, nil
}

func (r *UserSqliteRepository) GetByName(name string) (*model.User, error) {
	var userName string
	var userPassword string
	var userRole model.Role

	err := r.db.QueryRow(`
    SELECT
        U.NAME,
        U.PASSWORD_HASH,
        U.ROLE
    FROM USERS U
    WHERE U.NAME = ?;`, name).Scan(&userName, &userPassword, &userRole)

	if err != nil {
		return nil, err
	}

	return &model.User{
		Name:     userName,
		Password: userPassword,
		Role:     userRole,
	}, nil
}
