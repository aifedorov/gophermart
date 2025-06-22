package user

import (
	"errors"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidLogin       = errors.New("invalid login format")
	ErrInvalidPassword    = errors.New("invalid password format")
	ErrWeakPassword       = errors.New("password is too weak")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserInactive       = errors.New("user account is inactive")
)
