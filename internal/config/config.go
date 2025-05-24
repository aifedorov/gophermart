package config

import (
	"flag"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

const dotEnvFile = ".env"

type Config struct {
	ListenAddress        string `env:"RUN_ADDRESS" envDefault:":8080"`
	StorageDSN           string `env:"DATABASE_URI,required,notEmpty"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:":8081"`
	LogLevel             string `env:"LOG_LEVEL" envDefault:"info"`
}

func LoadConfig() (*Config, error) {
	var listenAddress, storageDSN, accrualSystemAddress, logLevel string

	flag.StringVar(&listenAddress, "a", "", "address and port to run server")
	flag.StringVar(&storageDSN, "d", "", "postgres connection string")
	flag.StringVar(&accrualSystemAddress, "r", "", "address and port to run accrual server")
	flag.StringVar(&logLevel, "l", "", "log level")
	flag.Parse()

	// Ignoring error because the `.env` file is not required.
	_ = godotenv.Load(dotEnvFile)

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	if listenAddress != "" {
		cfg.ListenAddress = listenAddress
	}
	if storageDSN != "" {
		cfg.StorageDSN = storageDSN
	}
	if accrualSystemAddress != "" {
		cfg.AccrualSystemAddress = accrualSystemAddress
	}
	if logLevel != "" {
		cfg.LogLevel = logLevel
	}

	return &cfg, nil
}
