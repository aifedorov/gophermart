package order

import (
	"context"

	"github.com/aifedorov/gophermart/internal/api"
	"github.com/aifedorov/gophermart/internal/domain/user"
)

type Service interface {
	CreateOrder(ctx context.Context, userID user.ID, number Number) error
	GetUserOrders(ctx context.Context, userID user.ID) ([]*api.OrderResponse, error)
}

type orderService struct {
	repo Repository
}

func NewOrderService(repo Repository) Service {
	return &orderService{
		repo: repo,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, userID user.ID, number Number) error {
	return nil
}

func (s *orderService) GetUserOrders(ctx context.Context, userID user.ID) ([]*api.OrderResponse, error) {
	return nil, nil
}
