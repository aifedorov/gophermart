package handlers

import (
	"github.com/aifedorov/gophermart/internal/api"
	"github.com/aifedorov/gophermart/internal/logger"
	"github.com/aifedorov/gophermart/internal/server/middleware/auth"
	"go.uber.org/zap"
	"net/http"

	"github.com/aifedorov/gophermart/internal/repository"
)

func NewGetOrdersHandler(repo repository.Repository) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "app/json")

		userID, err := auth.GetUserID(req)
		if err != nil {
			logger.Log.Info("user not authenticated", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		orders, err := repo.GetOrdersByUserID(userID)
		if err != nil {
			logger.Log.Error("failed to get orders", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusOK)
		if err := encodeResponse(rw, api.ToOrdersResponse(orders)); err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}
