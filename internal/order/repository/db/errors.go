package repository

import "errors"

var (
	ErrOrderAlreadyExists        = errors.New("order already exists")
	ErrOrderAddedByAnotherUser   = errors.New("order uploaded by another user")
	ErrWithdrawInsufficientFunds = errors.New("withdraw insufficient funds")
)
