package repository

import "errors"

var (
	ErrNotFound              = errors.New("user not found")
	ErrAlreadyExists         = errors.New("user already exists")
	ErrInvalidateCredentials = errors.New("login or password is invalid")
)

type Repository interface {
	CreateUser(login, password string) (User, error)
	GetUserByCredentials(login, password string) (User, error)
	CreateOrderByUserID(userID, orderNumber string) error
	GetOrdersByUserID(userID string) ([]Order, error)
}
