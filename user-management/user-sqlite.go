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

func (r *UserSqliteRepository) GetAll() ([]string, error) {
	var users []string

	rows, err := r.db.Query(`
    SELECT
        U.NAME
    FROM USERS U;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userName string

		err := rows.Scan(&userName)
		if err != nil {
			return nil, err
		}

		users = append(users, userName)
	}

	return users, nil
}

func (r *UserSqliteRepository) Create(user *model.User) error {
	_, err := r.db.Exec(`
    INSERT INTO USERS (NAME, PASSWORD_HASH, ROLE)
    VALUES (?, ?, ?);`,
		user.Name,
		user.Password,
		user.Role)

	return err
}

func (r *UserSqliteRepository) DeleteByName(name string) error {
	_, err := r.db.Exec(`
    DELETE FROM USERS
    WHERE NAME = ?;`, name)

	return err
}

func (r *UserSqliteRepository) Update(user *model.User) error {
	_, err := r.db.Exec(`
    UPDATE USERS
    SET PASSWORD_HASH = ?, ROLE = ?
    WHERE NAME = ?;`,
		user.Password,
		user.Role,
		user.Name)

	return err
}
