package repository

import "time"

type OrderStatus string

const (
	New        OrderStatus = "new"
	Processing OrderStatus = "processing"
	Processed  OrderStatus = "processed"
	Invalid    OrderStatus = "invalid"
)

type User struct {
	ID       string
	Login    string
	Password string
	Balance  float64
}

type Order struct {
	ID          string
	UserID      string
	Number      string
	Status      OrderStatus
	Amount      float64
	CreatedAt   time.Time
	ProcessedAt time.Time
}
