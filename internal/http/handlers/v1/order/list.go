package order

import (
	"context"
	"encoding/json"
	"net/http"

	authDomain "github.com/aifedorov/gophermart/internal/domain/auth"
	orderDomain "github.com/aifedorov/gophermart/internal/domain/order"
)

type ListHandler struct {
	orderService orderDomain.Service
	authService  authDomain.Service
}

func NewListHandler(orderService orderDomain.Service, authService authDomain.Service) *ListHandler {
	return &ListHandler{
		orderService: orderService,
		authService:  authService,
	}
}

func (h *ListHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := h.authService.GetUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	orders, err := h.orderService.GetUserOrders(context.TODO(), userID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(orders)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
