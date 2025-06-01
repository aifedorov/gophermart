package handlers

import (
	"errors"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/aifedorov/gophermart/internal/logger"
	"github.com/aifedorov/gophermart/internal/repository"
	"github.com/aifedorov/gophermart/pkg/validation"
)

func NewCreateOrdersHandler(repo repository.Repository) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "text/plain")

		order, err := io.ReadAll(req.Body)
		if err != nil {
			logger.Log.Info("failed to read request body", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if len(order) == 0 {
			logger.Log.Info("empty order number")
			http.Error(rw, "empty order number", http.StatusBadRequest)
			return
		}

		if !validation.IsValidOrderNumber(string(order)) {
			logger.Log.Info("invalid order number")
			http.Error(rw, "invalid order number", http.StatusUnprocessableEntity)
			return
		}

		err = repo.CreateOrder(string(order))
		if errors.Is(err, repository.ErrAlreadyExists) {
			logger.Log.Info("order already exists", zap.String("order", string(order)))
			rw.WriteHeader(http.StatusOK)
			return
		}
		if err != nil {
			logger.Log.Error("failed to store order", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusAccepted)
	}
}
