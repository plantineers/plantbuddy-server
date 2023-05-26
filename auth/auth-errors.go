package auth

import "errors"

var ErrWrongCredentials = errors.New("wrong credentials")
var ErrNoCredentials = errors.New("no credentials supplied")
var ErrInvalidAuthHeader = errors.New("invalid authorization header")
var ErrCannotDeleteRoot = errors.New("cannot delete root user")
var ErrUserAlreadyExists = errors.New("user already exists")
