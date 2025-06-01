package repository

import "errors"

var (
	ErrNotFound              = errors.New("user not found")
	ErrAlreadyExists         = errors.New("user already exists")
	ErrInvalidateCredentials = errors.New("login or password is invalid")
)

type Repository interface {
	StoreUser(login, password string) error
	FetchUser(login, password string) (User, error)
}
