package api

import "github.com/aifedorov/gophermart/internal/repository"

func ToOrdersResponse(orders []repository.Order) []OrderResponse {
	if len(orders) == 0 {
		return nil
	}

	respOrders := make([]OrderResponse, len(orders))
	for i, order := range orders {
		var accrual *float64
		if order.Amount > 0 {
			accrual = &order.Amount
		}

		respOrder := OrderResponse{
			Number:     order.Number,
			Status:     string(order.Status),
			Accrual:    accrual,
			UploadedAt: order.CreatedAt,
		}

		respOrders[i] = respOrder
	}
	return respOrders
}
