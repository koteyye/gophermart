package config

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	// Уровень логирования
	Level slog.Level `env:"LOG_LEVEL"`

	// Адрес запуска сервера
	RunAddress string `env:"RUN_ADDRESS"`

	// Строка подключения к БД
	DatabaseURI string `env:"DATABASE_URI"`
}

// Clone возвращает копию конфигурации
func (c *Config) Clone() *Config {
	c2 := *c
	return &c2
}

// Validate возвращает ошибку, если одно из полей конфигурации не валидно
func (c *Config) Validate() error {
	if c.RunAddress == "" {
		return errors.New("the run address must not be empty")
	}
	if c.DatabaseURI == "" {
		return errors.New("the database uri must not be empty")
	}
	return nil
}

// Parse парсит переменные окружения и устанавливает их в переданную конфигурацию
func Parse(c *Config) error {
	c2 := c.Clone()
	err := env.Parse(c)
	if err != nil {
		*c = *c2
		return fmt.Errorf("parsing env: %w", err)
	}
	return nil
}
