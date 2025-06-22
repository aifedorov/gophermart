package api

import (
	"github.com/aifedorov/gophermart/internal/domain/order"
)

func ToDomainOrdersResponse(orders []*order.Order) []OrderResponse {
	if len(orders) == 0 {
		return nil
	}

	respOrders := make([]OrderResponse, len(orders))
	for i, o := range orders {
		var accrual *float64
		if o.Status == order.StatusProcessed && o.Amount > 0 {
			accrual = &o.Amount
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
