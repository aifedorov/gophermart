package domain

import (
	"time"

	repository "github.com/aifedorov/gophermart/internal/order/repository/db"
)

func convertOrderToDomain(dbOrder repository.Order) Order {
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
		Status:      convertStatusToDomain(dbOrder.Status),
		Accrual:     dbOrder.Amount,
		CreatedAt:   createdAt,
		ProcessedAt: processedAt,
	}
}

func convertStatusToDomain(dbStatus repository.Orderstatus) Status {
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

func convertOrderToWithdrawalDomain(dbOrder repository.Order) (Withdrawal, error) {
	return Withdrawal{
		ID:          dbOrder.ID.String(),
		UserID:      dbOrder.UserID.String(),
		OrderNumber: dbOrder.Number,
		Sum:         dbOrder.Amount,
		ProcessedAt: dbOrder.ProcessedAt.Time,
	}, nil
}
