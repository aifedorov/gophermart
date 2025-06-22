package order

import (
	"fmt"

	"github.com/aifedorov/gophermart/internal/repository"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateOrder(userID, number string) (*Order, error) {
	if !isValidOrderNumber(number) {
		return nil, ErrInvalidOrderNumber
	}

	existingOrder, err := s.repo.GetOrderByNumber(number)
	if err == nil {
		if existingOrder.UserID == userID {
			return nil, ErrOrderAlreadyUploaded
		}
		return nil, ErrOrderUploadedByAnotherUser
	}

	repoOrder, err := s.repo.CreateOrder(userID, number)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return s.convertFromRepoOrder(repoOrder), nil
}

func (s *Service) GetUserOrders(userID string) ([]*Order, error) {
	repoOrders, err := s.repo.GetOrdersByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}

	orders := make([]*Order, 0, len(repoOrders))
	for _, repoOrder := range repoOrders {
		orders = append(orders, s.convertFromRepoOrder(repoOrder))
	}

	return orders, nil
}

func (s *Service) convertFromRepoOrder(repoOrder repository.Order) *Order {
	return &Order{
		ID:          repoOrder.ID,
		UserID:      repoOrder.UserID,
		Number:      repoOrder.Number,
		Status:      Status(repoOrder.Status),
		Amount:      repoOrder.Amount,
		CreatedAt:   repoOrder.CreatedAt,
		ProcessedAt: repoOrder.ProcessedAt,
	}
}

func (s *Service) ToOrderResponse(order *Order) Response {
	resp := Response{
		Number:     order.Number,
		Status:     string(order.Status),
		UploadedAt: order.CreatedAt,
	}

	if order.Status == StatusProcessed && order.Amount > 0 {
		resp.Amount = order.Amount
	}

	if !order.ProcessedAt.IsZero() {
		resp.ProcessedAt = order.ProcessedAt
	}

	return resp
}
