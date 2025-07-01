package handlers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/aifedorov/gophermart/internal/api"
	"github.com/aifedorov/gophermart/internal/domain/order"
	"github.com/aifedorov/gophermart/internal/logger"
	"github.com/aifedorov/gophermart/internal/server/middleware/auth"
)

func NewGetOrdersHandler(orderService *order.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		userID, err := auth.GetUserID(req)
		if err != nil {
			logger.Log.Info("user not authenticated", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		orders, err := orderService.GetUserOrders(userID)
		if err != nil {
			logger.Log.Error("failed to get orders", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if len(orders) == 0 {
			rw.WriteHeader(http.StatusNoContent)
			return
		}

		rw.WriteHeader(http.StatusOK)
		if err := encodeResponse(rw, api.ToOrdersResponse(orders)); err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}
