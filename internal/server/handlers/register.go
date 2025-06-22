package handlers

import (
	"errors"
	"github.com/aifedorov/gophermart/internal/config"
	"github.com/aifedorov/gophermart/internal/server/middleware/auth"
	"net/http"

	"go.uber.org/zap"

	"github.com/aifedorov/gophermart/internal/logger"
	"github.com/aifedorov/gophermart/internal/repository"
)

func NewRegisterHandler(cfg config.Config, repo repository.Repository) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		body, err := decodeRegister(req)
		if err != nil {
			logger.Log.Info("failed to decode request", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if !isValidCredentials(body.Login, body.Password) {
			logger.Log.Info("empty login or password")
			http.Error(rw, "empty login or password", http.StatusBadRequest)
			return
		}

		user, err := repo.CreateUser(body.Login, body.Password)
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

		auth.SetNewAuthCookies(user.ID, cfg.SecretKey, rw)
		rw.WriteHeader(http.StatusOK)
	}
}
