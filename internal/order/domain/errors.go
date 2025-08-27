package domain

import "errors"

var (
	ErrInvalidOrderNumber        = errors.New("invalid order number format")
	ErrOrderNotFound             = errors.New("order not found")
	ErrWithdrawNegativeAmount    = errors.New("withdraw amount should be positive")
	ErrWithdrawInsufficientFunds = errors.New("withdraw insufficient funds")
)
