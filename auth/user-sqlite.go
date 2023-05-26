// Author: Maximilian Floto
package auth

import (
	"database/sql"
	"errors"

	"github.com/plantineers/plantbuddy-server/db"
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

func (r *UserSqliteRepository) GetById(id int64) (*User, error) {
	var userId int64
	var userName string
	var userPassword string
	var userRole Role

	err := r.db.QueryRow(`
    SELECT
        U.ID,
        U.NAME,
        U.PASSWORD,
        U.ROLE
    FROM USERS U
    WHERE U.ID = ?;`, id).Scan(&userId, &userName, &userPassword, &userRole)

	if err != nil {
		return nil, err
	}

	return &User{
		Id:       userId,
		Name:     userName,
		Password: userPassword,
		Role:     userRole,
	}, nil
}

func (r *UserSqliteRepository) GetByName(name string) (*User, error) {
	var userId int64
	var userName string
	var userPassword string
	var userRole Role

	err := r.db.QueryRow(`
    SELECT
        U.ID,
        U.NAME,
        U.PASSWORD,
        U.ROLE
    FROM USERS U
    WHERE U.NAME = ?;`, name).Scan(&userId, &userName, &userPassword, &userRole)

	if err != nil {
		return nil, err
	}

	return &User{
		Id:       userId,
		Name:     userName,
		Password: userPassword,
		Role:     userRole,
	}, nil
}

func (r *UserSqliteRepository) GetAll() ([]int64, error) {
	var users []int64

	rows, err := r.db.Query(`
    SELECT
        U.ID
    FROM USERS U
    ORDER BY ID;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userId int64

		err := rows.Scan(&userId)
		if err != nil {
			return nil, err
		}

		users = append(users, userId)
	}

	return users, nil
}

func (r *UserSqliteRepository) Create(user *User) error {
	_, err := r.db.Exec(`
    INSERT INTO USERS (NAME, PASSWORD, ROLE)
    VALUES (?, ?, ?);`,
		user.Name,
		user.Password,
		user.Role)

	return err
}

func (r *UserSqliteRepository) DeleteById(id int64) error {
	_, err := r.db.Exec(`
    DELETE FROM USERS
    WHERE ID = ?;`, id)

	return err
}

func (r *UserSqliteRepository) Update(user *User) error {
	_, err := r.db.Exec(`
    UPDATE USERS
    SET PASSWORD = ?, ROLE = ?, NAME = ?
    WHERE ID = ?;`,
		user.Password,
		user.Role,
		user.Name,
		user.Id)

	return err
}
