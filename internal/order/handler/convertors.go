package handler

import (
	"github.com/aifedorov/gophermart/internal/order/domain"
	"github.com/shopspring/decimal"
)

func ToOrdersResponse(orders []domain.Order) []OrderResponse {
	if len(orders) == 0 {
		return nil
	}

	respOrders := make([]OrderResponse, len(orders))
	for i, o := range orders {
		var accrual *decimal.Decimal
		if o.Status == domain.StatusProcessed && o.Accrual.IsPositive() {
			accrual = &o.Accrual
		}

		respOrder := OrderResponse{
			Number:     o.Number,
			Status:     string(o.Status),
			Accrual:    accrual,
			UploadedAt: o.CreatedAt,
		}

		respOrders[i] = respOrder
	}
	return respOrders
}
