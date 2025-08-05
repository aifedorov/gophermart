package domain

import "errors"

var (
	ErrOrderAlreadyUploaded       = errors.New("order already uploaded by this user")
	ErrOrderUploadedByAnotherUser = errors.New("order uploaded by another user")
	ErrInvalidOrderNumber         = errors.New("invalid order number format")
	ErrOrderNotFound              = errors.New("order not found")
	ErrWithdrawNegativeAmount     = errors.New("withdraw amount should be positive")
	ErrWithdrawInsufficientFunds  = errors.New("withdraw insufficient funds")
)
