package handler

import (
	"errors"
	"github.com/aifedorov/gophermart/internal/pkg/config"
	"github.com/aifedorov/gophermart/internal/pkg/logger"
	"github.com/aifedorov/gophermart/internal/pkg/middleware"
	"github.com/aifedorov/gophermart/internal/user/domain"
	"github.com/aifedorov/gophermart/internal/user/repository"
	"net/http"

	"go.uber.org/zap"
)

func NewRegisterHandler(cfg config.Config, userService *domain.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		body, err := decodeRegister(req)
		if err != nil {
			logger.Log.Info("failed to decode request", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		userReq := repository.RegisterRequest{
			Login:    body.Login,
			Password: body.Password,
		}

		registeredUser, err := userService.Register(userReq)
		if errors.Is(err, domain.ErrEmptyCredentials) {
			logger.Log.Info("empty login or password")
			http.Error(rw, "empty login or password", http.StatusBadRequest)
			return
		}
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			logger.Log.Info("login already exists", zap.String("login", body.Login))
			http.Error(rw, "login already exists", http.StatusConflict)
			return
		}
		if err != nil {
			logger.Log.Error("failed to register user", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		middleware.SetNewAuthCookies(registeredUser.ID, cfg.SecretKey, rw)
		rw.WriteHeader(http.StatusOK)
	}
}
