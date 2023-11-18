package config

import (
	"errors"
	"flag"
	"log/slog"
)

type Config struct {
	// Уровень логирования
	Level slog.Level `env:"LOG_LEVEL"`

	// Адрес запуска сервера
	RunAddress string `env:"RUN_ADDRESS"`

	// Строка подключения к БД
	DatabaseURI string `env:"DATABASE_URI"`
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

func (c *Config) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.RunAddress, "a", "", "run address")
	fs.StringVar(&c.DatabaseURI, "d", "", "database uri")
	fs.TextVar(&c.Level, "v", slog.LevelInfo, "logging level")
}