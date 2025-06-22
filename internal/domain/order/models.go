package order

import "time"

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

type CreateOrderRequest struct {
	Number string `json:"number"`
}

type Response struct {
	Number      string    `json:"number"`
	Status      string    `json:"status"`
	Amount      float64   `json:"accrual,omitempty"`
	UploadedAt  time.Time `json:"uploaded_at"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
}
