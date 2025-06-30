package handlers

import (
	"net/http"

	"github.com/aifedorov/gophermart/internal/domain/user"
	"github.com/aifedorov/gophermart/internal/server/middleware/auth"
	"go.uber.org/zap"

	"github.com/aifedorov/gophermart/internal/logger"
)

func NewGetBalanceHandler(userService *user.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		userID, err := auth.GetUserID(req)
		if err != nil {
			logger.Log.Info("user not authenticated", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		balance, err := userService.GetUserBalance(userID)
		if err != nil {
			logger.Log.Error("failed to get balance", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusOK)
		if err := encodeJSONResponse(rw, balance); err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}
