package handlers

import (
	"github.com/aifedorov/gophermart/internal/app/config"
)

func newMockConfig() config.Config {
	return config.Config{
		ListenAddress:        "localhost:8080",
		StorageDSN:           "postgres://test",
		AccrualSystemAddress: "localhost:8081",
		LogLevel:             "debug",
		SecretKey:            "secret",
	}
}
