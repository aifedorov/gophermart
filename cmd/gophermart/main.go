package main

import (
	"context"
	"log"

	"github.com/aifedorov/gophermart/internal/client/accrual"
	orderDomain "github.com/aifedorov/gophermart/internal/order/domain"
	orderRepository "github.com/aifedorov/gophermart/internal/order/repository/db"
	"github.com/aifedorov/gophermart/internal/pkg/config"
	"github.com/aifedorov/gophermart/internal/pkg/logger"
	"github.com/aifedorov/gophermart/internal/pkg/posgre"
	"github.com/aifedorov/gophermart/internal/server"
	userDomain "github.com/aifedorov/gophermart/internal/user/domain"
	userRepository "github.com/aifedorov/gophermart/internal/user/repository/db"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = logger.Initialize(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = logger.Log.Sync()
	}()

	ctx := context.Background()
	db := posgre.NewPosgresRepository(ctx, cfg.StorageDSN)
	err = db.Open()
	defer db.Close()

	if err != nil {
		logger.Log.Fatal("failed to connect to database", zap.Error(err))
	}

	userRepo := userRepository.NewRepository(ctx, db.DBPool())
	userService := userDomain.NewService(userRepo)

	accrualClient := accrual.NewHTTPClient(cfg)
	defer func() {
		err := accrualClient.Close()
		if err != nil {
			logger.Log.Error("accrualclient: error closing http client", zap.Error(err))
		}
	}()

	orderRepo := orderRepository.NewRepository(ctx, db.DBPool())
	orderService := orderDomain.NewService(orderRepo)

	poller := orderDomain.NewPoller(ctx, orderRepo, accrualClient)
	checker := orderDomain.NewChecker(ctx, orderRepo, poller)
	go func() {
		err := checker.Run()
		if err != nil {
			logger.Log.Fatal("checker: error running checker", zap.Error(err))
		}
	}()

	s := server.NewServer(cfg, userService, orderService)
	if err := s.Run(); err != nil {
		logger.Log.Fatal("server: failed to run", zap.Error(err))
	}

	// TODO: Handle shoutdown
}
