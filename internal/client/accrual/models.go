package accrual

import "github.com/shopspring/decimal"

type Status string

const (
	StatusRegistered Status = "REGISTERED"
	StatusProcessing Status = "PROCESSING"
	StatusProcessed  Status = "PROCESSED"
	StatusInvalid    Status = "INVALID"
)

type OrderResponse struct {
	Number string           `json:"number"`
	Status Status           `json:"status"`
	Amount *decimal.Decimal `json:"amount,omitempty"`
}
