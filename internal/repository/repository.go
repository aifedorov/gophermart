package repository

import "errors"

var (
	ErrNotFound              = errors.New("user not found")
	ErrAlreadyExists         = errors.New("user already exists")
	ErrInvalidateCredentials = errors.New("login or password is invalid")
	ErrOrderNotFound         = errors.New("order not found")
)

// TODO: Split into domains.
type Repository interface {
	CreateUser(login, password string) (User, error)
	GetUserByCredentials(login, password string) (User, error)
	CreateOrder(userID, number string) (Order, error)
	GetOrdersByUserID(userID string) ([]Order, error)
	GetOrderByNumber(number string) (Order, error)
}
