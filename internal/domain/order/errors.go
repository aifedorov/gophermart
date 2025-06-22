package order

import "errors"

var (
	ErrOrderAlreadyUploaded       = errors.New("order already uploaded by this user")
	ErrOrderUploadedByAnotherUser = errors.New("order uploaded by another user")
	ErrInvalidOrderNumber         = errors.New("invalid order number format")
)
