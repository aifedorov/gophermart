package handler

import (
	"errors"
	"github.com/aifedorov/gophermart/internal/order/domain"
	"github.com/aifedorov/gophermart/internal/pkg/logger"
	"github.com/aifedorov/gophermart/internal/pkg/middleware"
	domain2 "github.com/aifedorov/gophermart/internal/user/domain"
	"net/http"

	"go.uber.org/zap"
)

func NewWithdrawHandler(userService *domain2.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		body, err := decodeWithdraw(req)
		if err != nil {
			logger.Log.Info("failed to decode request", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		userID, err := middleware.GetUserID(req)
		if err != nil {
			logger.Log.Info("user not authenticated", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if !domain.IsValidOrderNumber(body.Order) {
			logger.Log.Info("invalid order number", zap.String("order", body.Order))
			http.Error(rw, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}

		err = userService.Withdraw(userID, body.Order, body.Sum)
		if errors.Is(err, domain2.ErrWithdrawNegativeAmount) {
			logger.Log.Info("negative amount of money to withdraw")
			http.Error(rw, http.StatusText(http.StatusPaymentRequired), http.StatusPaymentRequired)
			return
		}
		if errors.Is(err, domain2.ErrWithdrawInsufficientFunds) {
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
