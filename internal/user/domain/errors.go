package domain

import "errors"

var (
	ErrInvalidCredentials        = errors.New("invalid credentials")
	ErrUserAlreadyExists         = errors.New("user already exists")
	ErrEmptyCredentials          = errors.New("empty login or password")
	ErrNotFound                  = errors.New("user not found")
	ErrAlreadyExists             = errors.New("user already exists")
	ErrInvalidateCredentials     = errors.New("login or password is invalid")
	ErrWithdrawNegativeAmount    = errors.New("withdraw amount should be positive")
	ErrWithdrawInsufficientFunds = errors.New("insufficient funds")
)
