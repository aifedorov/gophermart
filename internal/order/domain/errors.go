package domain

import "errors"

var (
	ErrOrderAlreadyUploaded       = errors.New("order already uploaded by this user")
	ErrOrderUploadedByAnotherUser = errors.New("order uploaded by another user")
	ErrInvalidOrderNumber         = errors.New("invalid order number format")
	ErrOrderNotFound              = errors.New("order not found")
	ErrAlreadyExists              = errors.New("order already exists")
)
