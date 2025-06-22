package handlers

import (
	"errors"
	"github.com/aifedorov/gophermart/internal/config"
	"github.com/aifedorov/gophermart/internal/domain/user"
	"github.com/aifedorov/gophermart/internal/server/middleware/auth"
	"net/http"

	"go.uber.org/zap"

	"github.com/aifedorov/gophermart/internal/logger"
)

func NewRegisterHandler(cfg config.Config, userService *user.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		body, err := decodeRegister(req)
		if err != nil {
			logger.Log.Info("failed to decode request", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		userReq := user.RegisterRequest{
			Login:    body.Login,
			Password: body.Password,
		}

		registeredUser, err := userService.Register(userReq)
		if errors.Is(err, user.ErrEmptyCredentials) {
			logger.Log.Info("empty login or password")
			http.Error(rw, "empty login or password", http.StatusBadRequest)
			return
		}
		if errors.Is(err, user.ErrUserAlreadyExists) {
			logger.Log.Info("login already exists", zap.String("login", body.Login))
			http.Error(rw, "login already exists", http.StatusConflict)
			return
		}
		if err != nil {
			logger.Log.Error("failed to register user", zap.Error(err))
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		auth.SetNewAuthCookies(registeredUser.ID, cfg.SecretKey, rw)
		rw.WriteHeader(http.StatusOK)
	}
}
