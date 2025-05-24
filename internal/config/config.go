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
}

func LoadConfig() (*Config, error) {
	var listenAddress string
	var storageDSN string
	var accrualSystemAddress string

	flag.StringVar(&listenAddress, "a", "", "address and port to run server")
	flag.StringVar(&storageDSN, "d", "", "postgres connection string")
	flag.StringVar(&accrualSystemAddress, "r", "", "address and port to run accrual server")
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

	return &cfg, nil
}
