package user_management

import "github.com/plantineers/plantbuddy-server/model"

type UserRepository interface {
	// GetByName returns a user by its name.
	GetByName(name string) (*model.User, error)
}
