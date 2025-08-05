package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrEmptyCredentials   = errors.New("empty login or password")
	ErrNotFound           = errors.New("user not found")
)
