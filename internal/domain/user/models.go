package user

import "time"

type User struct {
	ID       string
	Login    string
	Password string
	Balance  float64
}

type Balance struct {
	Current   float64
	Withdrawn float64
}

type Withdrawal struct {
	ID          string
	UserID      string
	OrderNumber string
	Sum         float64
	ProcessedAt time.Time
}

type RegisterRequest struct {
	Login    string
	Password string
}

type LoginRequest struct {
	Login    string
	Password string
}
