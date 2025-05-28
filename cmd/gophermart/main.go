package main

import (
	"log"

	"github.com/aifedorov/gophermart/internal/config"
	"github.com/aifedorov/gophermart/internal/logger"
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
}
