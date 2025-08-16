package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Status string

const (
	StatusNew        Status = "new"
	StatusProcessing Status = "processing"
	StatusProcessed  Status = "processed"
	StatusInvalid    Status = "invalid"
)

type CreateStatus int

const (
	CreateStatusSuccess CreateStatus = iota
	CreateStatusAlreadyUploaded
	CreateStatusUploadedByAnotherUser
	CreateStatusFailed
)

type Order struct {
	ID          string
	UserID      string
	Number      string
	Status      Status
	Accrual     decimal.Decimal
	CreatedAt   time.Time
	ProcessedAt time.Time
}

type Balance struct {
	Current   decimal.Decimal
	Withdrawn decimal.Decimal
}

type Withdrawal struct {
	ID          string
	UserID      string
	OrderNumber string
	Sum         decimal.Decimal
	ProcessedAt time.Time
}
