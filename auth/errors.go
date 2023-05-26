package auth

import "errors"

// Either the username or the password is wrong.
var ErrWrongCredentials = errors.New("wrong credentials!")

// No credentials were supplied.
var ErrNoCredentials = errors.New("no credentials supplied!")

// The user does not have the required permissions/role.
var ErrInsufficientPermissions = errors.New("insufficient permissions!")

// The authorization header is invalid.
var ErrInvalidAuthHeader = errors.New("invalid authorization header!")
