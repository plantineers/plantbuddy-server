// Author: Maximilian Floto
package auth

// User represents a user in the database and is used internally only.
type User struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
}

// SafeUser represents a user in the database and is used for the API.
// It does not contain the password hash.
type SafeUser struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Role Role   `json:"role"`
}

// Users represents a list of users.
type Users struct {
	Users []string `json:"users"`
}

type Role int8

const (
	Admin Role = iota
	Gardener
)
