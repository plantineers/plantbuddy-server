// Author: Maximilian Floto
package auth

type UserRepository interface {
	// GetByName returns a user by its name.
	GetById(id int64) (*User, error)
	GetByName(name string) (*User, error)
	GetAll() ([]int64, error)
	Create(user *User) error
	DeleteById(id int64) error
	Update(user *User) error
}
