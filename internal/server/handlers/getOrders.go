package handlers

import (
	"github.com/aifedorov/gophermart/internal/api"
	"github.com/aifedorov/gophermart/internal/logger"
	"go.uber.org/zap"
	"net/http"

	"github.com/aifedorov/gophermart/internal/repository"
)

func NewGetOrdersHandler(repo repository.Repository) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		orders, err := repo.GetOrders()
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
