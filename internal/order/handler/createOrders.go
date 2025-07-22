package handler

import (
	"errors"
	"github.com/aifedorov/gophermart/internal/order/domain"
	"github.com/aifedorov/gophermart/internal/pkg/logger"
	"github.com/aifedorov/gophermart/internal/pkg/middleware"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func NewCreateOrdersHandler(orderService *domain.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "text/plain")

		userID, err := middleware.GetUserID(req)
		if err != nil {
			logger.Log.Info("user not authenticated", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		orderNumber, err := io.ReadAll(req.Body)
		if err != nil {
			logger.Log.Info("failed to read request body", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if len(orderNumber) == 0 {
			logger.Log.Info("empty order number")
			http.Error(rw, "empty order number", http.StatusBadRequest)
			return
		}

		_, err = orderService.CreateOrder(userID, string(orderNumber))
		if errors.Is(err, domain.ErrInvalidOrderNumber) {
			logger.Log.Info("invalid order number", zap.String("order", string(orderNumber)))
			http.Error(rw, "invalid order number", http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(err, domain.ErrOrderAlreadyUploaded) {
			logger.Log.Info("order already uploaded by this user", zap.String("order", string(orderNumber)))
			rw.WriteHeader(http.StatusOK)
			return
		}
		if errors.Is(err, domain.ErrOrderUploadedByAnotherUser) {
			logger.Log.Info("order uploaded by another user", zap.String("order", string(orderNumber)))
			http.Error(rw, "order uploaded by another user", http.StatusConflict)
			return
		}
		if err != nil {
			logger.Log.Error("failed to create order", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusAccepted)
	}
}
