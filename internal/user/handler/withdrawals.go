package handler

import (
	"encoding/json"
	"github.com/aifedorov/gophermart/internal/pkg/logger"
	"github.com/aifedorov/gophermart/internal/pkg/middleware"
	"github.com/aifedorov/gophermart/internal/user/domain"
	"net/http"

	"go.uber.org/zap"
)

func NewWithdrawalsHandler(userService *domain.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		userID, err := middleware.GetUserID(req)
		if err != nil {
			logger.Log.Info("user not authenticated", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

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
				Sum:         withdrawal.Sum,
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
