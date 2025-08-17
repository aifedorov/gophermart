package handler

import (
	"net/http"

	"github.com/aifedorov/gophermart/internal/order/domain"
	"github.com/aifedorov/gophermart/internal/pkg/logger"
	"github.com/aifedorov/gophermart/internal/pkg/middleware"
	"go.uber.org/zap"
)

func NewBalanceHandler(orderService domain.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		userID, err := middleware.GetUserID(req)
		if err != nil {
			logger.Log.Info("user not authenticated", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		balance, err := orderService.GetUserBalance(userID)
		if err != nil {
			logger.Log.Error("failed to get balance", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusOK)
		response := BalanceResponse{
			Current:   float32(balance.Current.InexactFloat64()),
			Withdrawn: float32(balance.Withdrawn.InexactFloat64()),
		}
		if err := encodeJSONResponse(rw, response); err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}
