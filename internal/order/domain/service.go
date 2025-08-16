package domain

import (
	"errors"
	"fmt"

	"github.com/aifedorov/gophermart/internal/order/repository/db"
	"github.com/shopspring/decimal"
)

type Service interface {
	CreateOrder(userID, number string) (Order, CreateStatus, error)
	GetUserOrders(userID string) ([]Order, error)
	GetUserBalance(userID string) (Balance, error)
	Withdraw(userID, orderNumber string, amount decimal.Decimal) (Withdrawal, CreateStatus, error)
	GetWithdrawals(userID string) ([]Withdrawal, error)
}
type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateOrder(userID, number string) (Order, CreateStatus, error) {
	if !IsValidOrderNumber(number) {
		return Order{}, CreateStatusFailed, ErrInvalidOrderNumber
	}

	dbOrder, isCreated, err := s.repo.CreateTopUpOrder(userID, number)
	if err != nil {
		return Order{}, CreateStatusFailed, fmt.Errorf("orderservice: failed to create order: %w", err)
	}

	domainOrder := convertOrderToDomain(dbOrder)
	if !isCreated {
		if dbOrder.UserID.String() == userID {
			return domainOrder, CreateStatusAlreadyUploaded, nil
		}
		return domainOrder, CreateStatusUploadedByAnotherUser, nil
	}
	return domainOrder, CreateStatusSuccess, nil
}

func (s *service) GetUserOrders(userID string) ([]Order, error) {
	dbOrders, err := s.repo.GetOrdersByUserID(userID)
	if errors.Is(err, repository.ErrOrderAlreadyExists) {
		return nil, ErrOrderNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("orderservice: failed to get user orders: %w", err)
	}

	result := make([]Order, 0, len(dbOrders))
	for _, dbOrder := range dbOrders {
		domainOrder := convertOrderToDomain(dbOrder)
		result = append(result, domainOrder)
	}

	return result, nil
}

func (s *service) GetUserBalance(userID string) (Balance, error) {
	balance, err := s.repo.GetUserBalanceByUserID(userID)
	if err != nil {
		return Balance{}, fmt.Errorf("orderservice: failed to get user balance: %w", err)
	}

	withdrawn, err := s.repo.GetUserWithdrawByUserID(userID)
	if err != nil {
		return Balance{}, fmt.Errorf("orderservice: failed to get user withdrawn: %w", err)
	}

	return Balance{
		Current:   balance,
		Withdrawn: withdrawn,
	}, nil
}

func (s *service) Withdraw(userID, orderNumber string, amount decimal.Decimal) (Withdrawal, CreateStatus, error) {
	if !IsValidOrderNumber(orderNumber) {
		return Withdrawal{}, CreateStatusFailed, ErrInvalidOrderNumber
	}
	if !amount.IsPositive() {
		return Withdrawal{}, CreateStatusFailed, ErrWithdrawNegativeAmount
	}

	order, err := s.repo.CreateWithdrawalOrder(userID, orderNumber, amount)
	if errors.Is(err, repository.ErrOrderAlreadyExists) {
		return Withdrawal{}, CreateStatusAlreadyUploaded, nil
	}
	if errors.Is(err, repository.ErrWithdrawInsufficientFunds) {
		return Withdrawal{}, CreateStatusFailed, ErrWithdrawInsufficientFunds
	}
	if err != nil {
		return Withdrawal{}, CreateStatusFailed, fmt.Errorf("orderservice: failed to create order: %w", err)
	}
	return Withdrawal{
		ID:          order.ID.String(),
		UserID:      order.UserID.String(),
		OrderNumber: order.Number,
		Sum:         order.Amount,
		ProcessedAt: order.ProcessedAt.Time,
	}, CreateStatusSuccess, nil
}

func (s *service) GetWithdrawals(userID string) ([]Withdrawal, error) {
	dbOrders, err := s.repo.GetWithdrawalsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("orderservice: failed to get user orders: %w", err)
	}

	domainWithdrawals := make([]Withdrawal, len(dbOrders))
	for i, dbOrder := range dbOrders {
		domainWithdrawals[i], err = convertOrderToWithdrawalDomain(dbOrder)
		if err != nil {
			return nil, fmt.Errorf("orderservice: failed to convert order to withdrawal domain: %w", err)
		}
	}
	return domainWithdrawals, nil
}
