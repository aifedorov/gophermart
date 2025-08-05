package domain

import (
	"errors"
	"fmt"
	"github.com/aifedorov/gophermart/internal/order/repository/db"
	"github.com/aifedorov/gophermart/internal/pkg/logger"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"time"
)

type Service interface {
	CreateOrder(userID, number string) (Order, error)
	GetUserOrders(userID string) ([]Order, error)
	GetUserBalance(userID string) (Balance, error)
	Withdraw(userID, orderNumber string, amount decimal.Decimal) (Withdrawal, error)
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

func (s *service) CreateOrder(userID, number string) (Order, error) {
	if !IsValidOrderNumber(number) {
		logger.Log.Info("orderservice: invalid order number")
		return Order{}, ErrInvalidOrderNumber
	}

	dbOrder, err := s.repo.CreateTopUpOrder(userID, number)
	if errors.Is(err, repository.ErrOrderAlreadyExists) {
		logger.Log.Info("orderservice: order already exists", zap.Error(err))
		return Order{}, ErrOrderAlreadyUploaded
	}
	if errors.Is(err, repository.ErrOrderAddedByAnotherUser) {
		logger.Log.Info("orderservice: order already exists", zap.Error(err))
		return Order{}, ErrOrderUploadedByAnotherUser
	}
	if err != nil {
		logger.Log.Error("orderservice: error creating order", zap.Error(err))
		return Order{}, fmt.Errorf("orderservice: failed to create order: %w", err)
	}

	domainOrder := s.convertOrderToDomain(dbOrder)
	return domainOrder, nil
}

func (s *service) GetUserOrders(userID string) ([]Order, error) {
	dbOrders, err := s.repo.GetOrdersByUserID(userID)
	if errors.Is(err, repository.ErrOrderAlreadyExists) {
		return nil, ErrOrderNotFound
	}
	if err != nil {
		logger.Log.Error("orderservice: failed to get orders for user", zap.String("user_id", userID))
		return nil, fmt.Errorf("orderservice: failed to get user orders: %w", err)
	}

	result := make([]Order, 0, len(dbOrders))
	for _, dbOrder := range dbOrders {
		domainOrder := s.convertOrderToDomain(dbOrder)
		result = append(result, domainOrder)
	}

	return result, nil
}

func (s *service) GetUserBalance(userID string) (Balance, error) {
	balance, err := s.repo.GetUserBalanceByUserID(userID)
	if err != nil {
		logger.Log.Error("orderservice: failed to get user balance", zap.Error(err))
		return Balance{}, err
	}

	withdrawn, err := s.repo.GetUserWithdrawByUserID(userID)
	if err != nil {
		logger.Log.Error("orderservice: failed to get user withdrawn", zap.Error(err))
		return Balance{}, err
	}

	return Balance{
		Current:   balance,
		Withdrawn: withdrawn,
	}, nil
}

func (s *service) Withdraw(userID, orderNumber string, amount decimal.Decimal) (Withdrawal, error) {
	if !IsValidOrderNumber(orderNumber) {
		logger.Log.Info("orderservice: invalid order number", zap.String("amount", amount.String()))
		return Withdrawal{}, ErrInvalidOrderNumber
	}
	if !amount.IsPositive() {
		logger.Log.Info("orderservice: invalid amount", zap.String("amount", amount.String()))
		return Withdrawal{}, ErrWithdrawNegativeAmount
	}

	order, err := s.repo.CreateWithdrawalOrder(userID, orderNumber, amount)
	if errors.Is(err, repository.ErrOrderAlreadyExists) {
		logger.Log.Info("orderservice: order already exists", zap.Error(err))
		return Withdrawal{}, ErrOrderAlreadyUploaded
	}
	if errors.Is(err, repository.ErrWithdrawInsufficientFunds) {
		logger.Log.Info("orderservice: insufficient funds", zap.Error(err))
		return Withdrawal{}, ErrWithdrawInsufficientFunds
	}
	if err != nil {
		return Withdrawal{}, fmt.Errorf("orderservice: failed to create order: %w", err)
	}
	return Withdrawal{
		ID:          order.ID.String(),
		UserID:      order.UserID.String(),
		OrderNumber: order.Number,
		Sum:         order.Amount,
		ProcessedAt: order.ProcessedAt.Time,
	}, nil
}

func (s *service) GetWithdrawals(userID string) ([]Withdrawal, error) {
	dbOrders, err := s.repo.GetWithdrawalsByUserID(userID)
	if err != nil {
		logger.Log.Error("orderservice: failed to get orders for user", zap.String("user_id", userID))
		return nil, err
	}

	domainWithdrawals := make([]Withdrawal, len(dbOrders))
	for i, dbOrder := range dbOrders {
		domainWithdrawals[i], err = s.convertOrderToWithdrawalDomain(dbOrder)
		if err != nil {
			logger.Log.Error("orderservice: failed to convert order to withdrawal domain", zap.Error(err))
			return nil, err
		}
	}
	return domainWithdrawals, nil
}

func (s *service) convertOrderToDomain(dbOrder repository.Order) Order {
	var processedAt time.Time
	if dbOrder.ProcessedAt.Valid {
		processedAt = dbOrder.ProcessedAt.Time
	}

	var createdAt time.Time
	if dbOrder.CreatedAt.Valid {
		createdAt = dbOrder.CreatedAt.Time
	}

	return Order{
		ID:          dbOrder.ID.String(),
		UserID:      dbOrder.UserID.String(),
		Number:      dbOrder.Number,
		Status:      s.convertStatusToDomain(dbOrder.Status),
		Accrual:     dbOrder.Amount,
		CreatedAt:   createdAt,
		ProcessedAt: processedAt,
	}
}

func (s *service) convertStatusToDomain(dbStatus repository.Orderstatus) Status {
	switch dbStatus {
	case repository.OrderstatusNEW:
		return StatusNew
	case repository.OrderstatusPROCESSING:
		return StatusProcessing
	case repository.OrderstatusPROCESSED:
		return StatusProcessed
	case repository.OrderstatusINVALID:
		return StatusInvalid
	default:
		return StatusNew
	}
}

func (s *service) convertOrderToWithdrawalDomain(dbOrder repository.Order) (Withdrawal, error) {
	return Withdrawal{
		ID:          dbOrder.ID.String(),
		UserID:      dbOrder.UserID.String(),
		OrderNumber: dbOrder.Number,
		Sum:         dbOrder.Amount,
		ProcessedAt: dbOrder.ProcessedAt.Time,
	}, nil
}
