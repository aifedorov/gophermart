package order

import (
	"fmt"
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
	if !IsValidOrderNumber(number) {
		return nil, ErrInvalidOrderNumber
	}

	existingOrder, err := s.repo.GetOrderByNumber(number)
	if err == nil {
		if existingOrder.UserID == userID {
			return nil, ErrOrderAlreadyUploaded
		}
		return nil, ErrOrderUploadedByAnotherUser
	}

	order, err := s.repo.CreateOrder(userID, number)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return &order, nil
}

func (s *Service) GetUserOrders(userID string) ([]*Order, error) {
	orders, err := s.repo.GetOrdersByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}

	result := make([]*Order, 0, len(orders))
	for i := range orders {
		result = append(result, &orders[i])
	}

	return result, nil
}
