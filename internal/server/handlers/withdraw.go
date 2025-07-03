package handlers

import (
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/aifedorov/gophermart/internal/domain/order"
	"github.com/aifedorov/gophermart/internal/domain/user"
	"github.com/aifedorov/gophermart/internal/logger"
	"github.com/aifedorov/gophermart/internal/server/middleware/auth"
)

func NewWithdrawHandler(userService *user.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		body, err := decodeWithdraw(req)
		if err != nil {
			logger.Log.Info("failed to decode request", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		userID, err := auth.GetUserID(req)
		if err != nil {
			logger.Log.Info("user not authenticated", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if !order.IsValidOrderNumber(body.Order) {
			logger.Log.Info("invalid order number", zap.String("order", body.Order))
			http.Error(rw, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}

		err = userService.Withdraw(userID, body.Order, body.Sum)
		if errors.Is(err, user.ErrWithdrawNegativeAmount) {
			logger.Log.Info("negative amount of money to withdraw")
			http.Error(rw, http.StatusText(http.StatusPaymentRequired), http.StatusPaymentRequired)
			return
		}
		if errors.Is(err, user.ErrWithdrawInsufficientFunds) {
			logger.Log.Info("insufficient funds to withdraw")
			http.Error(rw, http.StatusText(http.StatusPaymentRequired), http.StatusPaymentRequired)
			return
		}
		if err != nil {
			logger.Log.Error("failed to withdraw money", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusOK)
	}
}
