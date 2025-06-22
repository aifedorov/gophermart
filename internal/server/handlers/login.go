package handlers

import (
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/aifedorov/gophermart/internal/config"
	"github.com/aifedorov/gophermart/internal/domain/user"
	"github.com/aifedorov/gophermart/internal/logger"
	"github.com/aifedorov/gophermart/internal/server/middleware/auth"
)

func NewLoginHandler(cfg config.Config, userService *user.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		body, err := decodeLogin(req)
		if err != nil {
			logger.Log.Info("failed to decode request", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		userReq := user.LoginRequest{
			Login:    body.Login,
			Password: body.Password,
		}

		authenticatedUser, err := userService.Login(userReq)
		if errors.Is(err, user.ErrEmptyCredentials) {
			logger.Log.Info("empty login or password")
			http.Error(rw, "empty login or password", http.StatusBadRequest)
			return
		}
		if errors.Is(err, user.ErrInvalidCredentials) {
			logger.Log.Info("invalid login or password")
			http.Error(rw, "invalid login or password", http.StatusUnauthorized)
			return
		}
		if err != nil {
			logger.Log.Error("failed to authenticate user", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		auth.SetNewAuthCookies(authenticatedUser.ID, cfg.SecretKey, rw)
		rw.WriteHeader(http.StatusOK)
	}
}
