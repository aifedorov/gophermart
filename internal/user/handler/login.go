package handler

import (
	"errors"
	"github.com/aifedorov/gophermart/internal/pkg/config"
	"github.com/aifedorov/gophermart/internal/pkg/logger"
	"github.com/aifedorov/gophermart/internal/pkg/middleware"
	"github.com/aifedorov/gophermart/internal/user/domain"
	"go.uber.org/zap"
	"net/http"
)

func NewLoginHandler(cfg config.Config, userService domain.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		body, err := decodeLogin(req)
		if err != nil {
			logger.Log.Info("failed to decode request", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		userReq := domain.LoginRequest{
			Login:    body.Login,
			Password: body.Password,
		}

		authenticatedUser, err := userService.Login(userReq)
		if errors.Is(err, domain.ErrEmptyCredentials) {
			logger.Log.Info("empty login or password")
			http.Error(rw, "empty login or password", http.StatusBadRequest)
			return
		}
		if errors.Is(err, domain.ErrInvalidCredentials) {
			logger.Log.Info("invalid login or password")
			http.Error(rw, "invalid login or password", http.StatusUnauthorized)
			return
		}
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			logger.Log.Info("user already exists")
			http.Error(rw, "user already exists", http.StatusBadRequest)
			return
		}
		if err != nil {
			logger.Log.Error("failed to authenticate user", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		middleware.SetNewAuthCookies(authenticatedUser.ID, cfg.SecretKey, rw)
		rw.WriteHeader(http.StatusOK)
	}
}
