package order

import (
	"context"

	"github.com/aifedorov/gophermart/internal/domain/user"
)

type Repository interface {
	Create(ctx context.Context, order *Order) error
	GetByUserID(ctx context.Context, userID user.ID) ([]*Order, error)
}
