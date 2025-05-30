package repository

import "errors"

var (
	ErrNotFound      = errors.New("user not found")
	ErrAlreadyExists = errors.New("user already exists")
)

type Repository interface {
	StoreUser(login, password string) error
	FetchUser(login string) (User, error)
}
