package config

import (
	"errors"
	"flag"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address string `env:"RUN_ADDRESS"`
	DSN string `env:"DATABASE_URI"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

type cliFlag struct {
	addr string
	dsn string
	accAddr string
}

func GetConfig() (*Config, error) {
	cliFlag := &cliFlag{}

	//Парсим флаги если они есть
	flag.StringVar(&cliFlag.addr, "a", "", "server address flag")
	flag.StringVar(&cliFlag.addr, "d", "", "dsn flag")
	flag.StringVar(&cliFlag.addr, "r", "", "accrual address flag")
	flag.Parse()

	//Парсим ENV
	var config Config
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	if config.Address == "" {
		config.Address = cliFlag.addr
	}
	if config.DSN == "" {
		config.DSN = cliFlag.dsn
	}
	if config.AccrualAddress == "" {
		config.AccrualAddress = cliFlag.accAddr
	}

	if config.Address == "" && config.AccrualAddress == "" && config.DSN == "" {
		return nil, errors.New("config can't be empty")
	}

	return &config, nil
}