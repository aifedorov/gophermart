package api

import "time"

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type OrderRequest struct {
	Number string `json:"number"`
}

type OrderResponse struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    *float64  `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}
