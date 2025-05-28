package handlers

import (
	"github.com/aifedorov/gophermart/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

func NewRegisterHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		body, err := decodeRequest(req)
		if err != nil {
			logger.Log.Error("failed to decode request", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if body.Login == "" || body.Password == "" {
			logger.Log.Info("empty login or password")
			http.Error(rw, "empty login or password", http.StatusBadRequest)
			return
		}

		// TODO: Check storage for existing login
		if body.Login == "loginExists" {
			logger.Log.Info("login already exists")
			http.Error(rw, "login already exists", http.StatusConflict)
			return
		}

		// TODO: authorization and return cookies

		rw.WriteHeader(http.StatusOK)
	}
}
