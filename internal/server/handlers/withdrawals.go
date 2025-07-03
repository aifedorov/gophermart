package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/aifedorov/gophermart/internal/api"
	"github.com/aifedorov/gophermart/internal/domain/user"
	"github.com/aifedorov/gophermart/internal/logger"
	"github.com/aifedorov/gophermart/internal/server/middleware/auth"
)

func NewWithdrawalsHandler(userService *user.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		userID, err := auth.GetUserID(req)
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

		withdrawalResponses := make([]api.WithdrawalResponse, len(withdrawals))
		for i, withdrawal := range withdrawals {
			withdrawalResponses[i] = api.WithdrawalResponse{
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
