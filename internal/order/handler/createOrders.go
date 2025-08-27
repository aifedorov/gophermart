package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/aifedorov/gophermart/internal/order/domain"
	"github.com/aifedorov/gophermart/internal/pkg/logger"
	"github.com/aifedorov/gophermart/internal/pkg/middleware"

	"go.uber.org/zap"
)

func NewCreateOrdersHandler(orderService domain.Service) http.HandlerFunc {
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

		_, status, err := orderService.CreateOrder(userID, string(orderNumber))
		if errors.Is(err, domain.ErrInvalidOrderNumber) {
			logger.Log.Info("invalid order number", zap.String("order", string(orderNumber)))
			http.Error(rw, "invalid order number", http.StatusUnprocessableEntity)
			return
		}

		switch status {
		case domain.CreateStatusSuccess:
			rw.WriteHeader(http.StatusAccepted)
			return
		case domain.CreateStatusAlreadyUploaded:
			rw.WriteHeader(http.StatusOK)
		case domain.CreateStatusUploadedByAnotherUser:
			http.Error(rw, "order uploaded by another user", http.StatusConflict)
		case domain.CreateStatusFailed:
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
