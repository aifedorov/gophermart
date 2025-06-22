package order

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusNew        Status = "NEW"
	StatusProcessing Status = "PROCESSING"
	StatusInvalid    Status = "INVALID"
	StatusProcessed  Status = "PROCESSED"
)

type ID string
type Number string

type Order struct {
	ID         ID
	Number     Number
	UserID     string
	Status     Status
	Accrual    float64
	UploadedAt time.Time
	UpdatedAt  time.Time
}

func NewOrder(number, userID string) *Order {
	now := time.Now()
	return &Order{
		ID:         ID(uuid.NewString()),
		Number:     Number(number),
		UserID:     userID,
		Status:     StatusNew,
		Accrual:    0,
		UploadedAt: now,
		UpdatedAt:  now,
	}
}

func (o *Order) UpdateStatus(status Status) {
	o.Status = status
	o.UpdatedAt = time.Now()
}

func (o *Order) BelongsToUser(userID string) bool {
	return o.UserID == userID
}
