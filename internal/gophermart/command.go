package gophermart

import (
	"context"
	"flag"
	"fmt"

	"log/slog"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/config"
	"github.com/sergeizaitcev/gophermart/internal/gophermart/server"
)

// Command определяет команду с набором флагов для gophermart.
type Command struct {
	fs     *flag.FlagSet
	actual *config.Config
}

// NewCommand возвращает новый экземпляр Command.
func NewCommand() *Command {
	cmd := &Command{
		fs: flag.NewFlagSet("gophermart", flag.ContinueOnError),
		actual: &config.Config{
			Level:         slog.LevelInfo,
			SecretKeyPath: "../../secret_key.txt",
		},
	}
	cmd.init()
	return cmd
}

func (cmd *Command) init() {
	// Отключение автоматического срабатывания Usage
	// в случае возникновения ошибки при парсинге флагов.
	cmd.fs.Usage = func() {}

	a := cmd.actual

	cmd.fs.StringVar(&a.RunAddress, "a", a.RunAddress, "run address")
	cmd.fs.StringVar(&a.DatabaseURI, "d", a.DatabaseURI, "database uri")
	cmd.fs.StringVar(&a.AccrualSystemAddress, "r", a.AccrualSystemAddress, "accrual system address")
	cmd.fs.StringVar(&a.SecretKeyPath, "s", a.SecretKeyPath, "secret key path")
	cmd.fs.TextVar(&a.Level, "v", a.Level, "logging level")
	cmd.fs.DurationVar(&a.TokenTTL, "t", a.TokenTTL, "token lifetime")
}

// Usage печатает формат использования команды.
func (cmd *Command) Usage() {
	fmt.Fprintf(cmd.fs.Output(), "Usage of %s:\n", cmd.fs.Name())
	cmd.fs.PrintDefaults()
}

// Parse парсит флаги командной строки.
func (cmd *Command) Parse(args []string) error {
	c := cmd.actual.Clone()

	err := cmd.fs.Parse(args)
	if err != nil {
		cmd.actual = c
		return fmt.Errorf("parsing flags: %w", err)
	}

	return nil
}

// Run запускает gophermart и блокируется до тех пор, пока не сработает
// контекст или функция не вернёт ошибку.
func (cmd *Command) Run(ctx context.Context) error {
	c := cmd.actual.Clone()

	err := initConfig(c)
	if err != nil {
		return fmt.Errorf("init config: %w", err)
	}

	srv := server.New(c)
	return srv.Run(ctx)
}

func initConfig(c *config.Config) error {
	err := config.Parse(c)
	if err != nil {
		return fmt.Errorf("parsing: %w", err)
	}

	err = c.Validate()
	if err != nil {
		return fmt.Errorf("validation: %w", err)
	}

	return nil
}
