package main

import (
	"github.com/aifedorov/gophermart/internal/config"
	"github.com/aifedorov/gophermart/internal/logger"
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

}
