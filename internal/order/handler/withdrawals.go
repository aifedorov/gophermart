package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aifedorov/gophermart/internal/order/domain"
	"github.com/aifedorov/gophermart/internal/pkg/logger"
	"github.com/aifedorov/gophermart/internal/pkg/middleware"
	"go.uber.org/zap"
)

func NewWithdrawalsHandler(userService domain.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		userID, _ := middleware.GetUserID(req)
		withdrawals, err := userService.GetWithdrawals(userID)
		if err != nil {
			logger.Log.Error("failed to get withdrawals", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if len(withdrawals) == 0 {
			rw.WriteHeader(http.StatusNoContent)
			return
		}

		withdrawalResponses := make([]WithdrawalResponse, len(withdrawals))
		for i, withdrawal := range withdrawals {
			withdrawalResponses[i] = WithdrawalResponse{
				Order:       withdrawal.OrderNumber,
				Sum:         float32(withdrawal.Sum.InexactFloat64()),
				ProcessedAt: withdrawal.ProcessedAt,
			}
		}

		rw.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(rw).Encode(withdrawalResponses); err != nil {
			logger.Log.Error("failed to encode withdrawals", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}
