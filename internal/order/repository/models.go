package repository

import (
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
	Accrual     float64
	CreatedAt   time.Time
	ProcessedAt time.Time
}
