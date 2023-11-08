package config

import (
	"errors"
	"fmt"

	"log/slog"

	"github.com/caarlos0/env/v10"
)

// Config определяет конфигурацию для gophermart.
type Config struct {
	// Уровень логирования.
	Level slog.Level `env:"LOG_LEVEL"`

	// Адресс запуска сервера.
	RunAddress string `env:"RUN_ADDRESS"`

	// Строка подключения к БД.
	DatabaseURI string `env:"DATABASE_URI"`

	// Адресс системы начисления.
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

// Clone возвращает копию конфигурации.
func (c *Config) Clone() *Config {
	c2 := *c
	return &c2
}

// Validate возвращает ошибку, если одно из полей конфигурации не валидно.
func (c *Config) Validate() error {
	if c.RunAddress == "" {
		return errors.New("the run address must be not empty")
	}
	if c.DatabaseURI == "" {
		return errors.New("the database uri must be not empty")
	}
	if c.AccrualSystemAddress == "" {
		return errors.New("the address of the accrual system should not be empty")
	}
	return nil
}

// Parse парсит переменные окружения и устанавливает их в переданную конфигурацию.
func Parse(c *Config) error {
	c2 := c.Clone()
	err := env.Parse(c2)
	if err != nil {
		return fmt.Errorf("parsing env: %w", err)
	}
	*c = *c2
	return nil
}
