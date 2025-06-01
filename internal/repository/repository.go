package repository

import "errors"

var (
	ErrNotFound              = errors.New("user not found")
	ErrAlreadyExists         = errors.New("user already exists")
	ErrInvalidateCredentials = errors.New("login or password is invalid")
)

type Repository interface {
	CreateUser(login, password string) error
	GetUserByCredentials(login, password string) (User, error)
	CreateOrder(orderNumber string) error
	GetOrders() ([]Order, error)
}
