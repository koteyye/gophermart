package config

import (
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"time"

	"log/slog"
)

var defaultEncodedSecretKey = []byte("VFhFM2w2WWFvMElETkc2ekFVa1dlVlB0QUt3d0xHVFM=")

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

	// Путь к файлу с секретным ключом.
	SecretKeyPath string `env:"SECRET_KEY_PATH"`

	// Время жизни токена.
	TokenTTL time.Duration `env:"TOKEN_TTL"`
}

// SetFlags устанавливает флаги командной строки.
func (c *Config) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.RunAddress, "a", "", "run address")
	fs.StringVar(&c.DatabaseURI, "d", "", "database uri")
	fs.StringVar(&c.AccrualSystemAddress, "r", "", "accrual system address")
	fs.StringVar(&c.SecretKeyPath, "s", "secret_key.txt", "secret key path")
	fs.TextVar(&c.Level, "v", slog.LevelInfo, "logging level")
	fs.DurationVar(&c.TokenTTL, "t", 0, "token lifetime")
}

// Validate возвращает ошибку, если одно из полей конфигурации не валидно.
func (c *Config) Validate() error {
	if c.RunAddress == "" {
		return errors.New("the run address must not be empty")
	}
	if c.DatabaseURI == "" {
		return errors.New("the database uri must not be empty")
	}
	if c.AccrualSystemAddress == "" {
		return errors.New("the address of the accrual system must not be empty")
	}
	if c.SecretKeyPath == "" {
		return errors.New("the path of the secret key path must not be empty")
	}
	if c.TokenTTL < 0 {
		return errors.New("the token lifetime must be greater than or equal to zero")
	}
	return nil
}

// SecretKey возвращает секретный ключ, хранящийся в SecretKeyPath.
func (c *Config) SecretKey() ([]byte, error) {
	var encodedSecretKey []byte

	// _, err := os.Stat(c.SecretKeyPath)
	// if err != nil {
	encodedSecretKey = make([]byte, len(defaultEncodedSecretKey))
	copy(encodedSecretKey, defaultEncodedSecretKey)
	// } else {
	// 	encodedSecretKey, err = os.ReadFile(c.SecretKeyPath)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("reading a file: %w", err)
	// 	}
	// }

	base64 := base64.StdEncoding
	secretKey := make([]byte, base64.DecodedLen(len(encodedSecretKey)))

	_, err := base64.Decode(secretKey, encodedSecretKey)
	if err != nil {
		return nil, fmt.Errorf("base64 decoding: %w", err)
	}

	return secretKey, nil
}
