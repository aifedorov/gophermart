package repository

import (
	"context"
	"errors"

	"github.com/aifedorov/gophermart/internal/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type Repository interface {
	GetOrderByNumber(number string) (Order, error)
	GetFirstOrderByStatus(status Orderstatus) (Order, error)
	UpdateOrderStatusByNumber(number string, status Orderstatus, amount *decimal.Decimal) error
	GetOrdersByUserID(userID string) ([]Order, error)
	CreateTopUpOrder(userID, orderNumber string) (Order, bool, error)
	CreateWithdrawalOrder(userID, orderNumber string, amount decimal.Decimal) (Order, error)
	GetWithdrawalsByUserID(userID string) ([]Order, error)
	GetUserBalanceByUserID(userID string) (decimal.Decimal, error)
	GetUserWithdrawByUserID(userID string) (decimal.Decimal, error)
}

type service struct {
	ctx     context.Context
	queries *Queries
	pgpool  *pgxpool.Pool
}

func NewRepository(ctx context.Context, pgpool *pgxpool.Pool) Repository {
	return &service{
		ctx:     ctx,
		queries: New(pgpool),
		pgpool:  pgpool,
	}
}

func (s *service) GetOrderByNumber(number string) (Order, error) {
	return s.queries.GetOrderByNumber(s.ctx, number)
}

func (s *service) GetOrdersByUserID(userID string) ([]Order, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return s.queries.GetTopUpOrdersByUserID(s.ctx, id)
}

func (s *service) GetFirstOrderByStatus(status Orderstatus) (Order, error) {
	return s.queries.GetOrderByStatus(s.ctx, status)
}

func (s *service) UpdateOrderStatusByNumber(number string, status Orderstatus, amount *decimal.Decimal) error {
	var amountValue decimal.Decimal
	if amount != nil {
		amountValue = *amount
	}
	return s.queries.UpdateOrderByNumber(s.ctx, UpdateOrderByNumberParams{
		Number: number,
		Status: status,
		Amount: amountValue,
	})
}

func (s *service) CreateTopUpOrder(userID, orderNumber string) (Order, bool, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return Order{}, false, err
	}

	newOrder, err := s.queries.CreateTopUpOrder(s.ctx, CreateTopUpOrderParams{
		UserID: id,
		Number: orderNumber,
		Amount: decimal.NewFromInt(0),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			existingOrder, err := s.queries.GetOrderByNumber(s.ctx, orderNumber)
			if err != nil {
				return Order{}, false, err
			}
			return existingOrder, false, nil
		}
		return Order{}, false, err
	}
	return newOrder, true, nil
}

func (s *service) CreateWithdrawalOrder(userID, orderNumber string, amount decimal.Decimal) (Order, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return Order{}, err
	}

	tx, err := s.pgpool.Begin(s.ctx)
	if err != nil {
		return Order{}, err
	}
	defer func() {
		err = tx.Rollback(s.ctx)
		if err != nil {
			logger.Log.Error("db: failed to rollback transaction", zap.Error(err))
		}
	}()

	qtx := s.queries.WithTx(tx)
	balance, err := qtx.GetUserBalanceByUserID(s.ctx, id)
	if err != nil {
		return Order{}, err
	}

	if balance.LessThan(amount) {
		return Order{}, ErrWithdrawInsufficientFunds
	}

	newWithdraw, err := s.queries.Withdrawal(
		s.ctx,
		WithdrawalParams{
			id,
			orderNumber,
			amount,
		},
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return Order{}, ErrOrderAlreadyExists
		}
		return Order{}, err
	}

	if err = tx.Commit(s.ctx); err != nil {
		return Order{}, err
	}

	return newWithdraw, err
}

func (s *service) GetWithdrawalsByUserID(userID string) ([]Order, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return s.queries.GetWithdrawalsByUserID(s.ctx, id)
}

func (s *service) GetUserBalanceByUserID(userID string) (decimal.Decimal, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return decimal.Decimal{}, err
	}
	return s.queries.GetUserBalanceByUserID(s.ctx, id)
}

func (s *service) GetUserWithdrawByUserID(userID string) (decimal.Decimal, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return decimal.Decimal{}, err
	}
	return s.queries.GetUserWithdrawByUserID(s.ctx, id)
}
