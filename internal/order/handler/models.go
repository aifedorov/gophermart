package handler

import (
	"time"

	"github.com/shopspring/decimal"
)

type OrderRequest struct {
	Number string `json:"number"`
}

type OrderResponse struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    *float32  `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type WithdrawRequest struct {
	Order string          `json:"order"`
	Sum   decimal.Decimal `json:"sum"`
}

type WithdrawalResponse struct {
	Order       string    `json:"order"`
	Sum         float32   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type BalanceResponse struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}
