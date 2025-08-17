package handler

import (
	"github.com/aifedorov/gophermart/internal/order/domain"
)

func ToOrdersResponse(orders []domain.Order) []OrderResponse {
	if len(orders) == 0 {
		return nil
	}

	respOrders := make([]OrderResponse, len(orders))
	for i, o := range orders {
		var accrual *float32
		if o.Status == domain.StatusProcessed && o.Accrual.IsPositive() {
			accrualFloat := float32(o.Accrual.InexactFloat64())
			accrual = &accrualFloat
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
