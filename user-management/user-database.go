// Author: Maximilian Floto
package user_management

import "github.com/plantineers/plantbuddy-server/model"

type UserRepository interface {
	// GetByName returns a user by its name.
	GetById(id int64) (*model.User, error)
	GetByName(name string) (*model.User, error)
	GetAll() ([]int64, error)
	Create(user *model.User) error
	DeleteById(id int64) error
	Update(user *model.User) error
}
