package main

import (
	"github.com/aifedorov/gophermart/internal/repository"
	"log"

	"github.com/aifedorov/gophermart/internal/config"
	"github.com/aifedorov/gophermart/internal/logger"
	"github.com/aifedorov/gophermart/internal/server"
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

	repo := repository.NewInMemoryStorage()
	s := server.NewServer(cfg, repo, repo)
	s.Run()
}
