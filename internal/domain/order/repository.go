package order

import "github.com/aifedorov/gophermart/internal/repository"

type Repository interface {
	CreateOrder(userID, number string) (repository.Order, error)
	GetOrdersByUserID(userID string) ([]repository.Order, error)
	GetOrderByNumber(number string) (repository.Order, error)
}
