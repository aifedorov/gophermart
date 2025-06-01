package handlers

import (
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/aifedorov/gophermart/internal/logger"
	"github.com/aifedorov/gophermart/internal/repository"
)

func NewRegisterHandler(repo repository.Repository) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		body, err := decodeRequest(req)
		if err != nil {
			logger.Log.Error("failed to decode request", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if !isValidCredentials(body.Login, body.Password) {
			logger.Log.Info("empty login or password")
			http.Error(rw, "empty login or password", http.StatusBadRequest)
			return
		}

		err = repo.StoreUser(body.Login, body.Password)
		if errors.Is(err, repository.ErrAlreadyExists) {
			logger.Log.Info("login already exists", zap.String("login", body.Login))
			http.Error(rw, "login already exists", http.StatusConflict)
			return
		}
		if err != nil {
			logger.Log.Error("failed to fetch user", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// TODO: Set cookies

		rw.WriteHeader(http.StatusOK)
	}
}
