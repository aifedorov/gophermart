package order

import (
	"github.com/google/uuid"
	"time"
)

type Status string

const (
	StatusNew        Status = "new"
	StatusProcessing Status = "processing"
	StatusProcessed  Status = "processed"
	StatusInvalid    Status = "invalid"
)

type Order struct {
	ID          string
	UserID      string
	Number      string
	Status      Status
	Amount      float64
	CreatedAt   time.Time
	ProcessedAt time.Time
}

type Withdrawal struct {
	ID          string
	UserID      string
	OrderNumber string
	Sum         float64
	CreatedAt   time.Time
	ProcessedAt *time.Time
}

func NewWithdrawal(userID, orderNumber string, sum float64) *Withdrawal {
	return &Withdrawal{
		ID:          uuid.New().String(),
		UserID:      userID,
		OrderNumber: orderNumber,
		Sum:         sum,
		CreatedAt:   time.Now(),
	}
}
