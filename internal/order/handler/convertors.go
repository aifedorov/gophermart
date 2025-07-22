package handler

import (
	"github.com/aifedorov/gophermart/internal/order/repository"
)

func ToOrdersResponse(orders []*repository.Order) []OrderResponse {
	if len(orders) == 0 {
		return nil
	}

	respOrders := make([]OrderResponse, len(orders))
	for i, o := range orders {
		var accrual *float64
		if o.Status == repository.StatusProcessed && o.Accrual > 0 {
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
