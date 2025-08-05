package main

import (
	"context"
	orderDomain "github.com/aifedorov/gophermart/internal/order/domain"
	orderRepository "github.com/aifedorov/gophermart/internal/order/repository/db"
	"github.com/aifedorov/gophermart/internal/pkg/config"
	"github.com/aifedorov/gophermart/internal/pkg/logger"
	"github.com/aifedorov/gophermart/internal/pkg/posgres"
	"github.com/aifedorov/gophermart/internal/server"
	userDomain "github.com/aifedorov/gophermart/internal/user/domain"
	userRepository "github.com/aifedorov/gophermart/internal/user/repository/db"
	"go.uber.org/zap"
	"log"
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
	db := posgres.NewPosgresRepository(ctx, cfg.StorageDSN)
	err = db.Open()
	defer db.Close()

	if err != nil {
		logger.Log.Fatal("failed to connect to database", zap.Error(err))
	}

	userRepo := userRepository.NewRepository(ctx, db.DBPool())
	userService := userDomain.NewService(userRepo)

	orderRepo := orderRepository.NewRepository(ctx, db.DBPool())
	orderService := orderDomain.NewService(orderRepo)

	s := server.NewServer(cfg, userService, orderService)
	if err := s.Run(); err != nil {
		logger.Log.Fatal("server: failed to run", zap.Error(err))
	}
}
